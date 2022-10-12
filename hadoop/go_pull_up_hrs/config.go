package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"time"
)

type HostConfig struct {
	HBaseDaemonScript string `yaml:"hbase_daemon_script"`
	Master            struct {
		HTTPPort int      `yaml:"http_port"`
		User     string   `yaml:"user"`
		Passwd   string   `yaml:"passwd"`
		SSHPort  string   `yaml:"ssh_port"`
		Targets  []string `yaml:"targets"`
	} `yaml:"master"`
	Slave struct {
		HTTPPort int      `yaml:"http_port"`
		User     string   `yaml:"user"`
		Passwd   string   `yaml:"passwd"`
		SSHPort  string   `yaml:"ssh_port"`
		Targets  []string `yaml:"targets"`
	} `yaml:"slave"`
}

type ProgramConfig struct {
	CheckOrPullUp struct {
		FeishuBot         string                `yaml:"feishu_bot"`
		CheckInterval     int64                 `yaml:"check_interval"`
		HbaseRegionServer map[string]HostConfig `yaml:"hbase_region_server"`
	} `yaml:"check_or_pull_up"`
}

func loadProgramConfig(path string) *ProgramConfig {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("read program config file error, %s", err)
		return nil
	}
	var out ProgramConfig
	if err := yaml.Unmarshal(data, &out); err != nil {
		log.Errorf("decode program config file error, %s", err)
		return nil
	}
	return &out
}

func init() {
	cfg = loadProgramConfig("config.yml")
	if cfg == nil {
		os.Exit(-1)
	}
	interval = time.Duration(cfg.CheckOrPullUp.CheckInterval) * time.Second
}
