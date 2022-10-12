package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	logging "github.com/ipfs/go-log/v2"
	"github.com/magicst0ne/alertmanager-webhook-feishu/feishu"
	"github.com/magicst0ne/alertmanager-webhook-feishu/model"
	"github.com/magicst0ne/alertmanager-webhook-feishu/utils"
	"github.com/tidwall/gjson"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"io"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	configFile = kingpin.Flag(
		"config.file",
		"configuration file path.",
	).Short('c').Default("config.yml").String()

	serverPort = kingpin.Flag(
		"web.listen-address",
		"Address to listen on",
	).Short('p').Default(":8086").String()

	sc = &SafeConfig{
		C: &Config{},
	}

	L = logging.Logger("<Feishu-AlertManager")
)

type BatchGetID struct {
	Code int `json:"code"`
	Data struct {
		UserList []struct {
			Mobile string `json:"mobile"`
			UserID string `json:"user_id,omitempty"`
		} `json:"user_list"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	level := map[string]logging.LogLevel{
		"debug": logging.LevelDebug,
		"info":  logging.LevelInfo,
		"error": logging.LevelError,
		"warn":  logging.LevelWarn,
	}
	setLevel, ok := level[strings.ToLower(os.Getenv("FEISHU_LOG_LEVEL"))]

	if ok {
		logging.SetAllLoggers(setLevel)
	} else {
		logging.SetAllLoggers(logging.LevelDebug)
	}

	gin.SetMode(gin.ReleaseMode)
	// load config  first time
	if err := sc.ReloadConfig(*configFile); err != nil {
		L.Fatalf("failed to load config file, %s", err)
	}
	r := gin.Default()
	r.POST("/webhook", func(c *gin.Context) {

		bodyRaw, err := io.ReadAll(c.Request.Body)
		if err != nil {
			L.Errorf("read request body error, %s", err)
		}

		L.Debugf(string(bodyRaw))
		var alertMsg model.AlertMessage
		err = json.Unmarshal([]byte(bodyRaw), &alertMsg)
		if err != nil {
			L.Errorf("failed to parse WebHookMessage, %s", err)
			c.JSON(200, gin.H{"ret": "-1", "msg": "invalid data"})
		} else {

			accessToken := c.Query("access_token")
			receiver := c.Query("receiver")

			receiverConfig, err := sc.GetConfigByName(receiver)
			if err != nil {
				c.JSON(200, gin.H{"ret": "-1", "msg": "receiver not exists"})
				L.Errorf("receiver(%s) does not exists", receiver)
				return
			}

			if accessToken != receiverConfig.AccessToken {
				c.JSON(200, gin.H{"ret": "-1", "msg": "invalid access_token"})
				L.Errorf("invalid access_token(%s)", accessToken)
				return
			}

			forward := alertMsg.Alerts.ForwardNotify()

			if forward != "" {
				receiverConfig, err = sc.GetConfigByName(forward)
				if err != nil {
					c.JSON(200, gin.H{"ret": "-1", "msg": "forward receiver not found"})
					L.Errorf("forward receiver(%s) does not exists", forward)
					return
				}
			}

			webhook, _ := feishu.NewFeishu(receiverConfig.Fsurl)
			webhookMessage := model.WebhookMessage{AlertMessage: alertMsg}
			webhookMessage.AlertHosts = make(map[string]string)
			token := sc.GetTenantAccessToken()
			mentions := token.GetUserIDByMobilesOrEmails(receiverConfig.Mentions.Mobiles, receiverConfig.Mentions.Emails)
			if mentions != nil {
				webhookMessage.OpenIDs = mentions
			}
			setPanel(&webhookMessage)
			err = webhook.Send(&webhookMessage)

			if err != nil {
				c.JSON(200, gin.H{"ret": "-1", "msg": "unknown error " + err.Error()})
				L.Errorf("unknown error, %s", err)
				return
			}
			c.JSON(200, gin.H{"ret": "0", "msg": "ok"})
		}
	})

	L.Fatal(r.Run(*serverPort))
}

func setPanel(message *model.WebhookMessage) {
	panel := message.Alerts.Panel()
	from, to := parseTimeRange(panel.TimeRange)
	if from == "" || to == "" {
		L.Error("time range parse failed")
		return
	}
	params := url.Values{}
	params.Add("from", from)
	params.Add("to", to)
	params.Add("panelId", panel.PanelID)
	params.Add("width", "1000")
	params.Add("height", "500")
	params.Add("tz", "Asia/Shanghai")

	u, err := url.Parse(sc.GetGrafana().GetURL())
	if err != nil {
		L.Errorf("url parse failed, %s", err)
		return
	}
	u.Path, err = url.JoinPath(u.Path, "render/d-solo", panel.Dashboard)
	if err != nil {
		L.Errorf("url join path failed, %s", err)
		return
	}

	u.RawQuery = params.Encode()
	if strings.HasPrefix(panel.SpecialParams, "var-instance") {
		u.RawQuery = params.Encode() + "&" + "var-instance=" + message.Alerts.GetInstance()
	} else {
		u.RawQuery = params.Encode() + "&" + panel.SpecialParams
	}
	filename := utils.RandStringBytesMaskImprSrc(10) + ".png"
	uri := u.String()
	L.Debugf("panel address: %s", uri)
	httpClient := resty.New()
	_, err = httpClient.SetTimeout(30 * time.Second).SetOutputDirectory(".images").R().SetOutput(filename).Get(uri)
	if err != nil {
		L.Errorf("download image failed, %s", err)
		return
	}
	absPath := path.Join(".images", filename)
	message.PanelURL = uploadImageToFeiShu(absPath)
}

func parseTimeRange(s string) (from, to string) {
	if s == "" {
		return "", ""
	}
	suffix := strings.ToLower(string(s[len(s)-1]))
	atoi, err := strconv.Atoi(string(s[0 : len(s)-1]))
	if err != nil {
		L.Errorf("string to int failed, %s", err)
		return "", ""
	}
	num := time.Duration(atoi)
	var val time.Duration
	switch suffix {
	// minutes
	case "m":
		val = num * time.Minute
	// hours
	case "h":
		val = num * time.Hour
	// days
	case "d":
		val = (num * time.Hour) * 24
	// weeks
	case "w":
		val = ((num * time.Hour) * 24) * 7
	default:
		val = num * time.Minute
	}
	now := time.Now()
	tFrom := now.UnixMilli() - val.Milliseconds()
	from = fmt.Sprintf("%d", tFrom)
	to = fmt.Sprintf("%d", now.UnixMilli())
	return
}

func uploadImageToFeiShu(filename string) (imageKey string) {
	if filename == "" {
		L.Error("filename cannot be empty")
		return ""
	}
	if sc == nil {
		L.Error("config undefined")
		return ""
	}
	const api = "https://open.feishu.cn/open-apis/im/v1/images"
	httpClient := resty.New()
	resp, err := httpClient.
		SetHeader("Content-Type", "multipart/form-data").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", sc.GetTenantAccessToken())).
		R().SetFile("image", filename).SetFormData(map[string]string{"image_type": "message"}).
		Post(api)
	if err != nil {
		L.Errorf("post feishu image upload api failed, %s", err)
		return ""
	}
	data := resp.Body()
	L.Debugf("resp: %s", string(data))
	ImageKey := gjson.GetBytes(data, "data.image_key").String()
	L.Debugf("image_key: %s", ImageKey)
	return ImageKey
}
