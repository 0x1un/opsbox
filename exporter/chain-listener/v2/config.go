package v2

import (
	logging "github.com/ipfs/go-log/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

var (
	Cfg Config
	log = logging.Logger("Chain-Listener")
)

type Node struct {
	URL    string `yaml:"url"`
	Name   string `yaml:"name"`
	Method string `yaml:"method"`
	Body   string `yaml:"body"`
}

type Config struct {
	Duration   time.Duration `yaml:"duration"`
	Blockchain []struct {
		PubGetFieldRule   map[string]string `yaml:"pub_get_field_rule"`
		PrivyGetFieldRule map[string]string `yaml:"privy_get_field_rule"`
		Name              string            `yaml:"name"`
		PubNode           []Node            `yaml:"pub_node"`
		PrivyNode         []Node            `yaml:"privy_node"`
	} `yaml:"blockchain"`
}

type Entry struct {
	Val float64
	Tag map[string]string
}

type Pool struct {
	Pub   Entry
	Privy Entry
}

func init() {
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("unable to read config.yml, %s", err)
	}
	if err := yaml.Unmarshal(data, &Cfg); err != nil {
		log.Fatalf("unable to parse config.yml, %s", err)
	}
}
