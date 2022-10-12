package model

import (
	prommodel "github.com/prometheus/common/model"
	"time"
)

// Data is the data passed to notification templates and webhook pushes.
//
// End-users should not be exposed to Go's type system, as this will confuse them and prevent
// simple things like simple equality checks to fail. Map everything to float64/string.
type AlertMessage struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   Alerts `json:"alerts"`

	GroupLabels       KV `json:"groupLabels"`
	CommonLabels      KV `json:"commonLabels"`
	CommonAnnotations KV `json:"commonAnnotations"`

	ExternalURL string `json:"externalURL"`
}

// Alert holds one alert for notification templates.
type Alert struct {
	Status       string    `json:"status"`
	Labels       KV        `json:"labels"`
	Annotations  KV        `json:"annotations"`
	StartsAt     time.Time `json:"startsAt"`
	EndsAt       time.Time `json:"endsAt"`
	GeneratorURL string    `json:"generatorURL"`
	Fingerprint  string    `json:"fingerprint"`
}

type GrafanaPanel struct {
	Dashboard     string `json:"dashboard"`
	PanelID       string `json:"panel_id"`
	TimeRange     string `json:"time_range"`
	SpecialParams string `json:"special_params"`
}

// Alerts is a list of Alert objects.
type Alerts []Alert

// Firing returns the subset of alerts that are firing.
func (as Alerts) Firing() []Alert {
	res := []Alert{}
	for _, a := range as {
		if a.Status == string(prommodel.AlertFiring) {
			res = append(res, a)
		}
	}
	return res
}

// Resolved returns the subset of alerts that are resolved.
func (as Alerts) Resolved() []Alert {
	res := []Alert{}
	for _, a := range as {
		if a.Status == string(prommodel.AlertResolved) {
			res = append(res, a)
		}
	}
	return res
}

// Resolved returns the subset of alerts that are resolved.
func (as Alerts) Severity() string {
	severity := "warning"
	for _, a := range as {
		if _, ok := a.Labels["severity"]; ok {
			if a.Labels["severity"] == "error" {
				severity = "error"
			} else if a.Labels["severity"] == "fatal" {
				severity = "fatal"
			}
		}
	}
	return severity
}

func (as Alerts) GetInstance() string {
	for _, a := range as {
		if v, ok := a.Labels["instance"]; ok {
			return v
		}
	}
	return ""
}

func (as Alerts) Panel() GrafanaPanel {
	panel := GrafanaPanel{}
	for _, a := range as {
		if v, ok := a.Labels["dashboard"]; ok {
			panel.Dashboard = v
		}
		if v, ok := a.Labels["panel_id"]; ok {
			panel.PanelID = v
		}
		if v, ok := a.Labels["time_range"]; ok {
			panel.TimeRange = v
		}
		if v, ok := a.Labels["special_params"]; ok {
			panel.SpecialParams = v
		}
	}
	return panel
}

func (as Alerts) SendNotify() string {
	notify := "false"
	for _, a := range as {
		if v, ok := a.Labels["send_notify"]; ok && v == "true" {
			return v
		}
	}
	return notify
}

func (as Alerts) ForwardNotify() string {
	for _, a := range as {
		if v, ok := a.Labels["forward_notify"]; ok {
			return v
		}
	}
	return ""
}
