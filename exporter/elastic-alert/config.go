package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var (
	CFG BuiltinConfig
)

type BuiltinConfig struct {
	Monitor struct {
		Promethues struct {
			PromMetric string `yaml:"prom_metric"`
			Port       int    `yaml:"port"`
		} `yaml:"promethues"`
	} `yaml:"monitor"`
	Elasticsearch struct {
		Address         []string `yaml:"address"`
		Username        string   `yaml:"username"`
		Password        string   `yaml:"password"`
		CertFingerPrint string   `yaml:"cert_finger_print"`
	} `yaml:"elasticsearch"`
	AlertRules []struct {
		Index string   `yaml:"index"`
		Keys  []string `yaml:"keys"`
		Dsl   string   `yaml:"dsl"`
	} `yaml:"alert_rules"`
}

func init() {
	path := "elastic-alert.yml"
	if data, err := ioutil.ReadFile(path); err != nil {
		log.Fatalf("unable to load config file, %s", err)
	} else {
		if err := yaml.Unmarshal(data, &CFG); err != nil {
			log.Fatalf("unable to decode config, %s", err)
		}
	}
}
