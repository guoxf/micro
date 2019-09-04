package micro

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	defaultClient               = "grpc"
	defaultClientRequestTimeout = "5s"
	defaultClientRetries        = 1
	defaultClientPoolSize       = 10
	defaultClientPoolTTL        = "1m"
	defaultRegisterTTL          = "30s"
	defaultRegisterInterval     = "10s"

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

//LoadFromEnv 如果未指定，从环境变量中获取
func (sc *ServiceConfig) LoadFromEnv() {
	if sc.Client == "" {
		sc.Client = os.Getenv("MICRO_CLIENT")
	}

	if sc.ClientRequestTimeout == "" {
		sc.ClientRequestTimeout = os.Getenv("MICRO_CLIENT_REQUEST_TIMEOUT")
	}

	if sc.ClientRetries <= 0 {
		retries := os.Getenv("MICRO_CLIENT_RETRIES")
		if retries != "" {
			r, err := strconv.ParseInt(retries, 10, 32)
			if err != nil {
				log.Println(err)
			} else {
				sc.ClientRetries = int(r)
			}
		}
	}

	if sc.ClientPoolSize <= 0 {
		str := os.Getenv("MICRO_CLIENT_POOL_SIZE")
		if str != "" {
			poolSize, err := strconv.ParseInt(str, 10, 32)
			if err != nil {
				log.Println(err)
			} else {
				sc.ClientPoolSize = int(poolSize)
			}
		}
	}

	if sc.ClientPoolTTL == "" {
		sc.ClientPoolTTL = os.Getenv("MICRO_CLIENT_POOL_TTL")
	}

	if sc.Server == "" {
		sc.Server = os.Getenv("MICRO_SERVER")
	}

	if sc.ServerAddress == "" {
		sc.ServerAddress = os.Getenv("MICRO_SERVER_ADDRESS")
	}

	if sc.ServerAdvertise == "" {
		sc.ServerAdvertise = os.Getenv("MICRO_SERVER_ADVERTISE")
	}

	if sc.ServerID == "" {
		sc.ServerID = os.Getenv("MICRO_SERVER_ID")
	}

	if len(sc.ServerMetaData) == 0 {
		metaData := os.Getenv("MICRO_SERVER_METADATA")
		if metaData != "" {
			sc.ServerMetaData = strings.Split(metaData, ",")
		}
	}

	if sc.ServerName == "" {
		sc.ServerName = os.Getenv("MICRO_SERVER_NAME")
	}

	if sc.ServerVersion == "" {
		sc.ServerVersion = os.Getenv("MICRO_SERVER_VERSION")
	}

	if sc.Broker == "" {
		sc.Broker = os.Getenv("MICRO_BROKER")
	}

	if len(sc.BrokerAddress) == 0 {
		brokerAddrs := os.Getenv("MICRO_BROKER_ADDRESS")
		if brokerAddrs != "" {
			sc.BrokerAddress = strings.Split(brokerAddrs, ",")
		}
	}

	if sc.Registry == "" {
		sc.Registry = os.Getenv("MICRO_REGISTRY")
	}

	if len(sc.RegistryAddress) == 0 {
		registryAddrs := os.Getenv("MICRO_REGISTRY_ADDRESS")
		if registryAddrs != "" {
			sc.RegistryAddress = strings.Split(registryAddrs, ",")
		}
	}

	if sc.Selector == "" {
		sc.Selector = os.Getenv("MICRO_SELECTOR")
	}

	if sc.Transport == "" {
		sc.Transport = os.Getenv("MICRO_TRANSPORT")
	}

	if len(sc.TransportAddress) == 0 {
		transportAddrs := os.Getenv("MICRO_TRANSPORT_ADDRESS")
		if transportAddrs != "" {
			sc.TransportAddress = strings.Split(transportAddrs, ",")
		}
	}
}

// LoadDefault 加载默认的配置
func (sc *ServiceConfig) LoadDefault() {
	if sc.ClientRequestTimeout == "" {
		sc.ClientRequestTimeout = defaultClientRequestTimeout
	}

	if sc.ClientPoolSize == 0 {
		sc.ClientPoolSize = defaultClientPoolSize
	}

	if sc.ClientPoolTTL == "" {
		sc.ClientPoolTTL = defaultClientPoolTTL
	}

	if sc.RegisterTTL == "" {
		sc.RegisterTTL = defaultRegisterTTL
	}

	if sc.RegisterInterval == "" {
		sc.RegisterInterval = defaultRegisterInterval
	}
}

// NewServiceConfig 配置顺序:文件->环境变量->默认值
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
	sc.LoadFromEnv()
	sc.LoadDefault()
	return &sc
}
