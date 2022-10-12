package main

import (
	"chain-listener/spider"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	logging "github.com/ipfs/go-log/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ybbus/jsonrpc/v3"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	httpClient = resty.New()
	log        = logging.Logger("<chain-listener>")
)

const (
	UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36"
)

var (
	BlockChain       = []string{"OMNI", "ETH", "BSC", "TRON", "BTC"}
	PrivateBlockFunc = []func() (float64, string){GetOMNILatestHeight, GetETHLatestHeight, GetBSCLatestHeight, GetTRONLatestHeight, GetBTCLatestHeight}
	MainnetBlockFunc = []func() (float64, string){GetMainnetOMNILatestHeight, GetMainnetETHLatestHeight, GetMainnetBSCLatestHeight, GetMainnetTRONLatestHeight, GetMainnetBTCLatestHeight}
)

func init() {
	httpClient.SetHeader("User-Agent", UA)
	httpClient.SetRetryCount(3).SetRetryWaitTime(10 * time.Second).SetRetryMaxWaitTime(20 * time.Second).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}

func main() {
	if len(BlockChain) != len(PrivateBlockFunc) || len(BlockChain) != len(MainnetBlockFunc) {
		return
	}
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		registry := prometheus.NewRegistry()
		for idx, block := range BlockChain {
			privyOpt := prometheus.GaugeOpts{
				Name: fmt.Sprintf("Private_%s", block),
				Help: "",
			}
			pValue, pFrom := PrivateBlockFunc[idx]()
			privyOpt.ConstLabels = prometheus.Labels{
				"from": pFrom,
			}
			pGauge := prometheus.NewGauge(privyOpt)
			pGauge.Set(pValue)
			if err := registry.Register(pGauge); err != nil {
				log.Errorf("register failed, %s, %s", privyOpt.Name, err)
			}
			// mainnet
			mainnetOpt := prometheus.GaugeOpts{
				Name: fmt.Sprintf("Mainnet_%s", block),
				Help: "",
			}
			mValue, mFrom := MainnetBlockFunc[idx]()
			mainnetOpt.ConstLabels = prometheus.Labels{
				"from": mFrom,
			}
			mGauge := prometheus.NewGauge(mainnetOpt)
			mGauge.Set(mValue)
			if err := registry.Register(mGauge); err != nil {
				log.Errorf("register failed, %s, %s", mainnetOpt.Name, err)
			}

			// diff
			diffOpt := prometheus.GaugeOpts{
				Name: fmt.Sprintf("%s_diff", block),
				Help: "",
			}
			diffOpt.ConstLabels = prometheus.Labels{
				"from": pFrom,
			}
			diffGauge := prometheus.NewGauge(diffOpt)

			diffGauge.Set(mValue - pValue)

			if err := registry.Register(diffGauge); err != nil {
				log.Errorf("register failed, %s, %s", diffOpt.Name, err)
			}
		}
		registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		registry.MustRegister(collectors.NewGoCollector())
		h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry})
		h.ServeHTTP(w, r)
	})
	log.Fatal(http.ListenAndServe(":9071", nil))
}

// Exporter

func GetETHLatestHeight() (float64, string) {
	return float64(getChainHeight("http://192.168.224.85:8545/", "eth_blockNumber")), "192.168.224.85:8545"
}

func GetBSCLatestHeight() (float64, string) {
	return float64(getChainHeight("http://192.168.224.78:8545/", "eth_blockNumber")), "192.168.224.78:8545"
}

func GetTRONLatestHeight() (float64, string) {
	return float64(getTronChainHeight("http://192.168.224.124:8090/wallet/getnodeinfo")), "192.168.224.124:8090"
}

func GetBTCLatestHeight() (float64, string) {
	return float64(getBTCChainHeight("http://btc:btc2021@192.168.224.124:8332")), "192.168.224.124:8332"
}

func GetOMNILatestHeight() (float64, string) {
	return float64(getOMNIChainHeight("http://192.168.224.122:8332/rest/chaininfo.json")), "192.168.224.122:8332"
}

func GetMainnetETHLatestHeight() (float64, string) {
	//return float64(getMainnetETHChainHeight("https://api.yitaifang.com/")), "yitaifang.com"
	return float64(getMainnetETHChainHeight("https://api.blockchair.com/ethereum/stats")), "blockchair.com"
}

func GetMainnetBSCLatestHeight() (float64, string) {
	val, tag := float64(spider.GetLastBlockFromBscscanDotCom()), "bscscan.com"
	if val < 0 {
		return float64(spider.GetLastBlockFromBSCXPLORERDotCom()), tag
	}
	return val, tag
}

func GetMainnetTRONLatestHeight() (float64, string) {
	return float64(getMainnetTronChainHeight("https://apilist.tronscanapi.com/api/system/homepage-bundle")), "tronscanapi.com"
}
func GetMainnetBTCLatestHeight() (float64, string) {
	return float64(getMainnetBTCChainHeight("https://api.blockchair.com/bitcoin/stats")), "blockchair.com"
}
func GetMainnetOMNILatestHeight() (float64, string) {
	val, tag := float64(getMainnetOMNIChainHeight("https://api.omniexplorer.info/v1/transaction/blocks")), "omniexplorer.info"
	if val < 0 {
		return getOMNIFromBlockChair(), tag
	}
	return val, tag
}

func getOMNIFromBlockChair() float64 {
	resp, err := httpClient.R().Get("https://api.blockchair.com/bitcoin/omni/stats")
	if err != nil {
		log.Errorf("get OMNI from blockchair failed, %s", err)
		return -1
	}
	var val BlockChairState

	if err := json.Unmarshal(resp.Body(), &val); err != nil {
		log.Errorf("blockchair mainnet body unmarshall failed, %s", err)
		return -1
	}

	return float64(val.Context.State)
}

func getMainnetBTCChainHeight(url string) int64 {
	if url == "" {
		return -1
	}
	resp, err := httpClient.R().Get(url)
	if err != nil {
		log.Errorf("get BTC from mainnet failed, %s", err)
		return -1
	}
	var val BlockChairState
	if err := json.Unmarshal(resp.Body(), &val); err != nil {
		log.Errorf("blockchair mainnet body unmarshall failed, %s", err)
		return -1
	}
	return val.Data.BestBlockHeight
}

func getMainnetOMNIChainHeight(url string) int64 {
	if url == "" {
		return -1
	}
	resp, err := httpClient.R().Get(url)
	if err != nil {
		log.Errorf("get omni from mainnet failed, %s", err)
		return -1
	}
	var val map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &val); err != nil {
		log.Errorf("omni-omniexplorer mainnet body unmarshall failed, %s", err)
		return -1
	}
	got, ok := val["latest"]
	if ok {
		switch t := got.(type) {
		case int64, int32, uint32, uint64, float64:
			return int64(got.(float64))
		default:
			log.Errorf("unexception latest-block type: %#v, exception: number", t)
			return -1
		}
	}
	log.Error("latest block not found in omniexplorer.com resp")
	return -1
}

func getMainnetETHChainHeight(url string) int64 {
	if url == "" {
		return -1
	}
	resp, err := httpClient.SetTimeout(time.Minute*3).
		SetHeader("Referer", "https://blockchair.com").
		SetHeader("Origin", "https://blockchair.com").
		R().Get(url)
	if err != nil {
		log.Errorf("get geth from mainnet failed, %s", err)
		return -1
	}
	var val BlockChairState

	if err := json.Unmarshal(resp.Body(), &val); err != nil {
		log.Errorf("geth-yitaifang mainnet body unmarshall failed, %s", err)
		return -1
	}
	return int64(val.Data.BestBlockHeight)
}

func getMainnetTronChainHeight(url string) int64 {
	resp, err := httpClient.
		R().Get(url)
	if err != nil {
		log.Errorf("get tron from mainnet failed, %s", err)
		return -1
	}
	var val TronMainnet
	if len(resp.Body()) == 0 {
		log.Errorf("the tron response from the mainnet is empty")
		return -1
	}
	if err := json.Unmarshal(resp.Body(), &val); err != nil {
		log.Errorf("tron mainnet body unmarshall failed, %s", err)
		return -1
	}
	return int64(val.Tps.Data.BlockHeight)
}

func getOMNIChainHeight(url string) int64 {
	if url == "" {
		return -1
	}
	resp, err := httpClient.R().Get(url)
	if err != nil {
		log.Errorf("get omni from rest api failed, %s", err)
		return -1
	}
	var val map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &val); err != nil {
		log.Errorf("omni response unmarshall failed, %s", err)
		return -1
	}
	got, ok := val["blocks"]
	if ok {
		switch t := got.(type) {
		case int64, int32, uint32, uint64, float64:
			return int64(got.(float64))
		default:
			log.Errorf("unexception blocks type: %#v, exception: number", t)
		}
	}
	log.Error("blocks not found in unmarshalled obj")
	return -1
}

func getBTCChainHeight(url string) int64 {
	var val map[string]interface{}
	err := getChainHeightToVal(url, "getblockchaininfo", &val)
	if err != nil {
		log.Error(err)
		return -1
	}
	got, ok := val["blocks"]
	if ok {
		switch t := got.(type) {
		case int64, int32, uint32, uint64, float64:
			return int64(got.(float64))
		default:
			log.Errorf("unexception blocks type: %#v, exception: number", t)
		}
	}
	log.Error("blocks not found in unmarshalled obj")
	return -1
}

func getTronChainHeight(url string) int64 {
	if url == "" {
		return -1
	}

	resp, err := httpClient.R().Get(url)
	if err != nil {
		log.Errorf("get tron request failed, %s", err)
		return -1
	}
	var obj map[string]interface{}
	if len(resp.Body()) == 0 {
		log.Errorf("get tron response body is empty")
		return -1
	}
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		log.Errorf("unmarshal tron response body failed, %s", err)
		return -1
	}
	got, ok := obj["beginSyncNum"]
	if ok {
		switch t := got.(type) {
		case int64, uint64, int, int32, float64:
			return int64(got.(float64))
		default:
			log.Errorf("unexception beginSyncNum type: %#v, exception: number", t)
			return -1
		}

	}
	log.Errorf("beginSyncNum not found in unmarshalled obj")
	return -1
}
func getChainHeightToVal(url string, method string, val interface{}) error {
	_rpc := jsonrpc.NewClient(url)
	err := _rpc.CallFor(context.Background(), val, method, []string{})
	if err != nil {
		return fmt.Errorf("call the jsonrpc failed, %s", err)
	}
	return nil
}

func getChainHeight(url string, method string) int64 {
	_rpc := jsonrpc.NewClient(url)
	var val string
	err := _rpc.CallFor(context.Background(), &val, method, []string{})
	if err != nil {
		log.Errorf("call the jsonrpc failed, %s", err)
		return -1
	}
	return int64(hex2int(val))
}

func hex2int(str string) uint64 {
	cleaned := strings.Replace(str, "0x", "", -1)
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return result
}
