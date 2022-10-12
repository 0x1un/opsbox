package spider

import (
	"bytes"
	"crypto/tls"
	"github.com/antchfx/htmlquery"
	"github.com/go-resty/resty/v2"
	logging "github.com/ipfs/go-log/v2"
	"strconv"
	"time"
)

var (
	httpClient = resty.New()
	log        = logging.Logger("<spider>")
)

const (
	UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36"
)

func init() {
	httpClient.SetRetryCount(3).SetRetryWaitTime(10 * time.Second).SetRetryMaxWaitTime(20 * time.Second).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}

func GetLastBlockFromBSCXPLORERDotCom() int64 {
	url := "https://bscxplorer.com"

	resp, err := httpClient.R().
		SetHeader("User-Agent", UA).
		Get(url)
	if err != nil {
		log.Errorf("get bscxplorer.com page failed, %s", err)
		return -1
	}
	nodes, err := htmlquery.Parse(bytes.NewReader(resp.Body()))
	if err != nil {
		log.Errorf("parse html failed, %s", err)
		return -1
	}
	node, err := htmlquery.Query(nodes, "//*[@id=\"wrap\"]/div/div/div/div/div[2]/div[3]/div[4]/p[2]/text()")
	if err != nil {
		log.Errorf("query html failed, %s", err)
		return -1
	}
	got, err := strconv.Atoi(node.Data)
	if err != nil {
		log.Errorf("string convert to number failed, %s", err)
		return -1
	}
	return int64(got)
}

func GetLastBlockFromBscscanDotCom() int64 {
	url := "https://bscscan.com"
	resp, err := httpClient.R().
		SetHeader("User-Agent", UA).
		Get(url)
	if err != nil {
		log.Errorf("get bscscan.com page failed, %s", err)
		return -1
	}
	nodes, err := htmlquery.Parse(bytes.NewReader(resp.Body()))
	if err != nil {
		log.Errorf("parse html failed, %s", err)
		return -1
	}
	node, err := htmlquery.Query(nodes, "//*[@id=\"lastblock\"]/text()")
	if err != nil {
		log.Errorf("query html failed, %s", err)
		return -1
	}
	if node == nil {
		log.Errorf("query data is empty")
		return -1
	}
	got, err := strconv.Atoi(node.Data)
	if err != nil {
		log.Errorf("string convert to number failed, %s", err)
		return -1
	}
	return int64(got)
}
