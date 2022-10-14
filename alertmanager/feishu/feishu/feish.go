package feishu

import (
	"bytes"
	"github.com/0x1un/opsbox/alertmanager/feishu/model"
	"github.com/0x1un/opsbox/alertmanager/feishu/tmpl"
	"text/template"
)

type Feishu struct {
	webhook  string
	sdk      *Sdk
	tpl      *template.Template
	alertTpl *template.Template
}

func NewFeishu(furl string) (*Feishu, error) {

	// template
	tpl, err := tmpl.GetEmbedTemplate("default.tmpl")
	alertTpl, err := tmpl.GetEmbedTemplate("default_alert.tmpl")

	if err != nil {
		return nil, err
	}

	return &Feishu{
		webhook:  furl,
		sdk:      NewSDK("", ""),
		tpl:      tpl,
		alertTpl: alertTpl,
	}, nil
}

func (b Feishu) Send(alerts *model.WebhookMessage) error {

	// prepare data
	err := b.preprocessAlerts(alerts)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = b.tpl.Execute(&buf, alerts)
	if err != nil {
		return err
	}
	return b.sdk.WebhookV2(b.webhook, &buf)
}

func (b Feishu) preprocessAlerts(webhookMsg *model.WebhookMessage) error {

	// preprocess using alert template
	webhookMsg.Severity = webhookMsg.Alerts.Severity()
	webhookMsg.SendNotify = webhookMsg.Alerts.SendNotify()

	n := 0
	for _, alert := range webhookMsg.Alerts.Firing() {
		webhookMsg.FiringAlerts = append(webhookMsg.FiringAlerts, alert)

		if _, ok1 := alert.Labels["hostname"]; ok1 {
			if _, ok2 := webhookMsg.AlertHosts[alert.Labels["hostname"]]; ok2 {
				if _, ok3 := alert.Labels["instance"]; ok3 {
					webhookMsg.AlertHosts[alert.Labels["hostname"]] = alert.Labels["instance"]
				}
			}
		}

		n++
	}
	webhookMsg.FiringNum = n

	for _, alert := range webhookMsg.Alerts.Resolved() {
		webhookMsg.ResolvedAlerts = append(webhookMsg.ResolvedAlerts, alert)
	}

	return nil
}
