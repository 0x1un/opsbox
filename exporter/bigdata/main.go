package main

import (
	"encoding/json"
	"fmt"
	"github.com/0x1un/opsbox/exporter/bigdata/config"
	"github.com/gobeam/stringy"
	logging "github.com/ipfs/go-log/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	log = logging.Logger("<bigdata_exporter>")
	cfg config.MetricConfig
)

func init() {
	_cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config file, %s", err)
	}
	cfg = *_cfg
}

type (
	Bean     map[string]any
	JMXBeans map[string][]Bean
)

func (b Bean) Get(metricName string) any {
	if value, ok := b[metricName]; ok {
		return value
	}
	return nil
}

const tableMetric = `Namespace_(\w+)_table_(\w+)_metric_(\w+)`

var reTable = regexp.MustCompile(tableMetric)

type hrsTables struct {
	Namespace  string
	TableName  string
	MetricName string
	Value      any
}

func (b Bean) GetTable() []hrsTables {
	out := make([]hrsTables, 0)
	for k, v := range b {
		if !reTable.MatchString(k) {
			continue
		}
		matched := reTable.FindStringSubmatch(k)
		if len(matched) != 4 {
			continue
		}

		out = append(out, hrsTables{
			Namespace:  matched[1],
			TableName:  matched[2],
			MetricName: matched[3],
			Value:      v,
		})
	}
	return out
}

func (j JMXBeans) Get(name string) Bean {
	for _, bean := range j["beans"] {
		if _name, ok := bean["name"]; ok {
			if _name == name {
				return bean
			}

			// DataNodeActivity 后面跟的是变量
			if strings.HasPrefix(name, "Hadoop:service=DataNode,name=DataNodeActivity") &&
				strings.HasPrefix(_name.(string), name) {
				return bean
			}
		}
	}
	return nil
}

func RetrieveMetricFromTarget(target string) (JMXBeans, error) {
	resp, err := http.Get("http://" + target + "/jmx")
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var out JMXBeans
	return out, json.Unmarshal(data, &out)
}

func metricHandler(w http.ResponseWriter, r *http.Request) {
	moduleName := r.URL.Query().Get("module")
	if moduleName == "" {
		http.Error(w, "module name is missing", http.StatusBadRequest)
		return
	}
	params := r.URL.Query()
	target := params.Get("target")
	if target == "" {
		http.Error(w, "target parameter is missing", http.StatusBadRequest)
		return
	}
	info, err := RetrieveMetricFromTarget(target)
	if err != nil {
		http.Error(w, "retrieve metric failed: "+target+" "+err.Error(), http.StatusBadRequest)
		return
	}
	metrics, ok := cfg.Modules[moduleName]
	if !ok {
		http.Error(w, moduleName+" module not defined", http.StatusBadRequest)
		return
	}
	if moduleName == "hbase_rs_tbs" {
		scrapeMultiMetric(info, w, r)
	} else {
		scrapeSingleMetric(info, metrics.ScrapeMetrics, w, r)
	}
}

func scrapeMultiMetric(info JMXBeans, w http.ResponseWriter, r *http.Request) {
	name := "Hadoop:service=HBase,name=RegionServer,sub=Tables"
	registry := prometheus.NewRegistry()
	for _, line := range info.Get(name).GetTable() {
		opts := prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_%s", "hadoop_hbase_rs_tbs", line.MetricName),
			Help: "",
		}
		opts.ConstLabels = prometheus.Labels{
			"data_table": fmt.Sprintf("%s.%s", line.Namespace, line.TableName),
		}
		gauge := prometheus.NewGauge(opts)
		switch line.Value.(type) {
		case float64, int64:
			gauge.Set(line.Value.(float64))
		}
		err := registry.Register(gauge)
		if err != nil {
			log.Errorf("register %s failed, %s", opts.Name+"..."+opts.ConstLabels["data_table"], err)
		}
	}
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func scrapeSingleMetric(info JMXBeans, metrics []string, w http.ResponseWriter, r *http.Request) {
	registry := prometheus.NewRegistry()
	for _, v := range metrics {
		name, metricName, helpMsg, TagName := parseScrapeMetric(v)
		if name == "" {
			continue
		}
		var (
			strVal         = "na"
			numVal float64 = -1
		)
		bean := info.Get(name)
		var value any

		opts := prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_%s", TagName, stringy.New(metricName).SnakeCase("?", "").ToLower()),
			Help: helpMsg,
		}

		value = bean.Get(metricName)

		switch value.(type) {
		case int64, int32, int, float64:
			numVal = value.(float64)
		case string:
			strVal = value.(string)
		default:
			continue
		}
		if strVal != "na" {
			if metricName == "tag.NumOpenConnectionsPerUser" {
				var _m map[string]int
				_ = json.Unmarshal([]byte(strVal), &_m)
				if len(_m) == 0 {
					continue
				}
				o := make(map[string]string, 0)
				for _k, _v := range _m {
					o[_k] = strconv.Itoa(_v)
				}
				opts.ConstLabels = o
			}
		}
		gauge := prometheus.NewGauge(opts)
		if numVal != -1 {
			gauge.Set(numVal)
		}
		registry.MustRegister(gauge)
	}
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

// name metric_name help_msg tag
func parseScrapeMetric(line string) (string, string, string, string) {
	sp := strings.Split(line, " ")
	if len(sp) == 4 {
		return sp[0], sp[1], sp[2], sp[3]
	}
	return "", "", "", ""
}

func main() {
	http.HandleFunc("/metrics", metricHandler)
	log.Fatal(http.ListenAndServe(":7070", nil))
}
