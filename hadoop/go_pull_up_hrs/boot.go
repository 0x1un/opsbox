package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

var (
	requests = resty.New()
)

func sendToFeishuBot(content string) {
	tpl := []byte(fmt.Sprintf(`{
        "msg_type": "interactive",
        "card": {
            "config": {
                "wide_screen_mode": true,
                "enable_forward": true
            },
            "elements": [{
                "tag": "div",
                "text": {
                    "content": "%s",
                    "tag": "lark_md"
                }
            }],
            "header": {
                "title": {
                    "content": "HRegionServer 监测服务",
                    "tag": "plain_text"
                }
            }
        }
    }`, content))
	resp, err := requests.R().SetHeader("Content-Type", "application/json").SetBody(tpl).Post(cfg.CheckOrPullUp.FeishuBot)
	if err != nil {
		log.Errorf("send feishu msg failed, %s", err)
	}
	log.Debug(string(resp.Body()))
}
