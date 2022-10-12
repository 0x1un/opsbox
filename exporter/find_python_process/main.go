package main

import (
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/v3/process"
	"net/http"
	"strings"
)

var (
	log = logging.Logger("<python-discover>")
)

type PythonProcessSeed struct {
	Cmdline string
}

func getPythonProcess() []PythonProcessSeed {
	out := make([]PythonProcessSeed, 0)
	pids, err := process.Pids()
	if err != nil {
		log.Errorf("unable to get pid, %s", err)
	}
	for _, pid := range pids {
		proc, err := process.NewProcess(pid)
		if err != nil {
			//log.Errorf("failed to describe %d, %s", pid, err)
			continue
		}
		name, err := proc.Name()
		if err != nil {
			log.Error(err)
			continue
		}
		if strings.ToLower(name) == "python" || strings.ToLower(name) == "python.exe" {
			cmdline, err := proc.Cmdline()
			if err != nil {
				log.Errorf("unable to get %s cmdline, %s", name, err)
				continue
			}
			out = append(out, PythonProcessSeed{
				fmt.Sprintf("%d_%s", pid, getScriptLocation(cmdline)),
			})
		}
	}
	return out
}

func getScriptLocation(s string) string {
	if s == "" {
		return ""
	}
	parts := strings.Split(s, " ")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}

func main() {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		registry := prometheus.NewRegistry()
		registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		registry.MustRegister(collectors.NewGoCollector())

		processes := getPythonProcess()
		for _, seed := range processes {
			opt := prometheus.GaugeOpts{Name: "pps_up"}
			opt.ConstLabels = map[string]string{"cmdline": seed.Cmdline}
			gauge := prometheus.NewGauge(opt)
			gauge.Set(1)
			registry.MustRegister(gauge)
		}

		if len(processes) == 0 {
			opt := prometheus.GaugeOpts{Name: "pps_up"}
			opt.ConstLabels = map[string]string{"cmdline": ""}
			gauge := prometheus.NewGauge(opt)
			gauge.Set(0)
			registry.MustRegister(gauge)
		}
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}).ServeHTTP(w, r)
	})
	log.Fatal(http.ListenAndServe(":9073", nil))
}
