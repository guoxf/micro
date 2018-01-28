package micro

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	defaultClient               = "grpc"
	defaultClientRequestTimeout = "5s"
	defaultClientRetries        = 1
	defaultClientPoolSize       = 10
	defaultClientPoolTTL        = "1m"

	defaultServer        = "grpc"
	defaultServerVersion = "1.0.0"
	defaultServerAddress = "0.0.0.0:0"

	defaultBroker    = "http"
	defaultRegistry  = "consul"
	defaultSelector  = "cache"
	defaultTransport = "grpc"
)

type ServiceConfig struct {
	Client               string   `yaml:"client"`
	ClientRequestTimeout string   `yaml:"clientRequestTimeout"`
	ClientRetries        int      `yaml:"clientRetries"`
	ClientPoolSize       int      `yaml:"clientPoolSize"`
	ClientPoolTTL        string   `yaml:"clientPoolTTL"`
	Server               string   `yaml:"server"`
	ServerName           string   `yaml:"serverName"`
	ServerVersion        string   `yaml:"serverVersion"`
	ServerID             string   `yaml:"serverID"`
	ServerAddress        string   `yaml:"serverAddress"`
	ServerAdvertise      string   `yaml:"serverAdvertise"`
	ServerMetaData       []string `yaml:"serverMetaData"`
	Broker               string   `yaml:"broker"`
	BrokerAddress        []string `yaml:"brokerAddress"`
	Registry             string   `yaml:"registry"`
	RegistryAddress      []string `yaml:"registryAddress"`
	Selector             string   `yaml:"selector"`
	Transport            string   `yaml:"transport"`
	TransportAddress     []string `yaml:"transportAddress"`
	RegisterTTL          string   `yaml:"registerTTL"`
	RegisterInterval     string   `yaml:"registerInterval"`
}

func NewServiceConfig(configPath string) *ServiceConfig {
	if configPath == "" {
		panic("config path is empty!")
	}

	var sc ServiceConfig
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(b, &sc); err != nil {
		panic(err)
	}
	return &sc
}
