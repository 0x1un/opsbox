package main

import (
	v2 "chain-listener/v2"
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	log = logging.Logger("main")
)

func main() {
	ticker := time.NewTicker(v2.Cfg.Duration * time.Second)
	stop := make(chan bool)
	go func() {
		for {
			if v2.SafeMap.Count() == 0 {
				v2.Loop()
			}
			select {
			case <-ticker.C:
				v2.Loop()
			case s := <-stop:
				if s {
					ticker.Stop()
					return
				}
			}
		}
	}()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		SetProm().ServeHTTP(w, r)
	})
	if err := http.ListenAndServe(":9071", nil); err != nil {
		stop <- true
		log.Error(err)
	}
}

func SetProm() http.Handler {
	registry := prometheus.NewRegistry()
	for chain := range v2.SafeMap.Items() {
		pool, ok := v2.SafeMap.Get(chain)
		if !ok {
			continue
		}
		// private
		privyOpt := prometheus.GaugeOpts{Name: fmt.Sprintf("Private_%s", chain), Help: ""}
		pValue := pool.Privy.Val
		privyOpt.ConstLabels = pool.Privy.Tag
		pGauge := prometheus.NewGauge(privyOpt)
		pGauge.Set(pValue)
		if err := registry.Register(pGauge); err != nil {
			log.Errorf("register failed, %s, %s", privyOpt.Name, err)
		}

		// public
		mainnetOpt := prometheus.GaugeOpts{Name: fmt.Sprintf("Mainnet_%s", chain), Help: ""}
		mValue := pool.Pub.Val
		mainnetOpt.ConstLabels = pool.Pub.Tag
		mGauge := prometheus.NewGauge(mainnetOpt)
		mGauge.Set(mValue)
		if err := registry.Register(mGauge); err != nil {
			log.Errorf("register failed, %s, %s", mainnetOpt.Name, err)
		}

		// diff

		diffOpt := prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_diff", chain),
			Help: "",
		}
		diffOpt.ConstLabels = privyOpt.ConstLabels

		diffGauge := prometheus.NewGauge(diffOpt)

		diffGauge.Set(mValue - pValue)

		if err := registry.Register(diffGauge); err != nil {
			log.Errorf("register failed, %s, %s", diffOpt.Name, err)
		}
	}
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry})
	return h
}
