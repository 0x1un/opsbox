package v2

import (
	"context"
	"github.com/0x1un/opsbox/exporter/chain-listener/spider"
	"github.com/go-resty/resty/v2"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/tidwall/gjson"
	"github.com/ybbus/jsonrpc/v3"
	"strconv"
	"strings"
)

var (
	httpClient = resty.New()
)

var (
	SafeMap = cmap.New[*Pool]()
)

const (
	UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36"
)

var (
	funcBox = map[string]func() int64{
		"spider.GetLastBlockFromBscscanDotCom":    spider.GetLastBlockFromBscscanDotCom,
		"spider.GetLastBlockFromBSCXPLORERDotCom": spider.GetLastBlockFromBSCXPLORERDotCom,
	}
)

func init() {
	httpClient.SetHeader("User-Agent", UA)
}

func Loop() {

	for _, blockchain := range Cfg.Blockchain {
		log.Infof("try to get %s", blockchain.Name)
		pubEntry := Entry{}
		privyEntry := Entry{}
		pubEntry.Val, pubEntry.Tag = getField(blockchain.PubNode, blockchain.PubGetFieldRule)
		privyEntry.Val, privyEntry.Tag = getField(blockchain.PrivyNode, blockchain.PrivyGetFieldRule)

		SafeMap.Set(blockchain.Name, &Pool{
			Pub:   pubEntry,
			Privy: privyEntry,
		})
	}
}

func getField(nodes []Node, rules map[string]string) (float64, map[string]string) {
	if len(nodes) == 0 || len(rules) == 0 {
		log.Error("not enough information, argument \"nodes\" or \"rules\" is required")
		return -1, map[string]string{"from": "unknown"}
	}

	var out float64
	outTags := make(map[string]string, 0)
	for _, node := range nodes {
		// for get
		tags := map[string]string{
			"from": node.Name,
		}
		if strings.ToLower(node.Method) == "get" {
			resp, err := httpClient.R().Get(node.URL)
			if err != nil {
				log.Errorf("unable to get %s, %s", node.Name, err)
				out = -1
				outTags["from"] = node.Name
				continue
			}
			return gjson.GetBytes(resp.Body(), rules[node.Name]).Float(), tags
		}
		// for json-rpc
		if strings.ToLower(node.Method) == "jsonrpc" {
			return float64(getChainHeight(node.URL, rules[node.Name])), tags
		}

		// for spider
		if strings.HasPrefix(node.Method, "spider") {
			return float64(funcBox[node.Method]()), tags
		}

		// for post

		if strings.ToLower(node.Method) == "post" {
			resp, err := httpClient.R().SetBody(node.Body).Post(node.URL)
			if err != nil {
				log.Errorf("unable to post %s, %s", node.Name, err)
				out = -1
				outTags["from"] = node.Name
				continue
			}
			return gjson.GetBytes(resp.Body(), rules[node.Name]).Float(), tags
		}
	}
	return out, outTags
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
