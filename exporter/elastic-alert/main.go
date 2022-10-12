package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	logging "github.com/ipfs/go-log/v2"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	log = logging.Logger("<elastic-alert>")
)

type Connector struct {
	es *elasticsearch.Client
}

func NewConnector() *Connector {
	es, err := NewESConnector(
		CFG.Elasticsearch.Address,
		CFG.Elasticsearch.Username,
		CFG.Elasticsearch.Password,
		CFG.Elasticsearch.CertFingerPrint,
	)
	if err != nil {
		log.Errorf("init es connector failed, %s", err)
		return nil
	}
	return &Connector{es: es}
}

func (con *Connector) Search(index, dsl string) (*gjson.Result, error) {
	result, err := con.es.Search(
		con.es.Search.WithContext(context.Background()),
		con.es.Search.WithIndex(index),
		con.es.Search.WithBody(bytes.NewReader([]byte(dsl))),
		con.es.Search.WithTrackTotalHits(true),
		con.es.Search.WithPretty())
	if err != nil {
		return nil, err
	}
	if result.IsError() {
		return nil, err
	}
	defer result.Body.Close()
	data, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	p := gjson.ParseBytes(data)
	return &p, nil
}

func main() {
	conn := NewConnector()
	if conn == nil {
		return
	}
	http.HandleFunc(CFG.Monitor.Promethues.PromMetric, func(w http.ResponseWriter, r *http.Request) {
		registry := prometheus.NewRegistry()
		for _, rule := range CFG.AlertRules {
			parsed, err := conn.Search(rule.Index, rule.Dsl)
			if err != nil {
				log.Error(err)
				continue
			}
			fields := parsed.Get("hits.hits.#.fields")
			ids := parsed.Get("hits.hits.#._id").Array()
			if fields.IsArray() {
				for idx, field := range fields.Array() {
					logOpt := prometheus.GaugeOpts{
						Name: strings.Replace(strings.Replace(fmt.Sprintf("log_%s", rule.Index), "*", "X", -1), "-", "_", -1),
					}
					logOpt.ConstLabels = prometheus.Labels{}
					logOpt.ConstLabels["_id"] = ids[idx].Raw
					for _, key := range rule.Keys {
						logOpt.ConstLabels[key] = fmt.Sprintf("%s", field.Get(key+"|@flatten|0"))
					}
					gauge := prometheus.NewGauge(logOpt)
					gauge.SetToCurrentTime()
					if err := registry.Register(gauge); err != nil {
						log.Errorf("register failed, %s, %s", err, logOpt.ConstLabels)
					}
				}
			}
		}
		registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		registry.MustRegister(collectors.NewGoCollector())
		h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry})
		h.ServeHTTP(w, r)
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", CFG.Monitor.Promethues.Port), nil))

}

func NewESConnector(addr []string, username, password string, certFp string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 30,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
		CertificateFingerprint: certFp,
		Addresses:              addr,
		Username:               username,
		Password:               password,

		RetryOnStatus: []int{429, 502, 503, 504},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	info, err := es.Info()
	if err != nil {
		return nil, err
	}
	if info.IsError() {
		return nil, errors.New(info.String())
	}
	return es, nil
}
