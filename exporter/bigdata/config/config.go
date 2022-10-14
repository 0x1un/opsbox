package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Metrics struct {
	ScrapeMetrics []string `yaml:"scrape_metrics"`
}

type MetricConfig struct {
	Modules map[string]Metrics `yaml:"modules"`
}
type MetricConfig2 struct {
	Modules struct {
		Namenode struct {
			ScrapeMetrics []string `yaml:"scrape_metrics"`
		} `yaml:"namenode"`
		Yarn struct {
			ScrapeMetrics []string `yaml:"scrape_metrics"`
		} `yaml:"yarn"`
		Hbase struct {
			ScrapeMetrics []string `yaml:"scrape_metrics"`
		} `yaml:"hbase"`
	} `yaml:"modules"`
}

func LoadConfig() (*MetricConfig, error) {
	data, err := os.ReadFile("config/config.yml")
	if err != nil {
		return nil, err
	}
	var mc MetricConfig
	return &mc, yaml.Unmarshal(data, &mc)
}
