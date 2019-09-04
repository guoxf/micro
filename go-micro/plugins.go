package micro

import (
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/broker/http"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/registry/mdns"
	"github.com/micro/go-micro/selector"
	"github.com/micro/go-micro/selector/cache"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/transport"
	"github.com/micro/go-plugins/registry/zookeeper"

	thttp "github.com/micro/go-micro/transport/http"
	"github.com/micro/go-plugins/broker/kafka"

	clientgrpc "github.com/micro/go-plugins/client/grpc"
	servergrpc "github.com/micro/go-plugins/server/grpc"
	transportgrpc "github.com/micro/go-plugins/transport/grpc"
)

var (
	DefaultBrokers = map[string]func(...broker.Option) broker.Broker{
		"http":  http.NewBroker,
		"kafka": kafka.NewBroker,
	}

	DefaultClients = map[string]func(...client.Option) client.Client{
		"rpc":  client.NewClient,
		"grpc": clientgrpc.NewClient,
	}

	DefaultRegistries = map[string]func(...registry.Option) registry.Registry{
		"consul":    consul.NewRegistry,
		"mdns":      mdns.NewRegistry,
		"zookeeper": zookeeper.NewRegistry,
	}

	DefaultSelectors = map[string]func(...selector.Option) selector.Selector{
		"default": selector.NewSelector,
		"cache":   cache.NewSelector,
	}

	DefaultServers = map[string]func(...server.Option) server.Server{
		"rpc":  server.NewServer,
		"grpc": servergrpc.NewServer,
	}

	DefaultTransports = map[string]func(...transport.Option) transport.Transport{
		"http": thttp.NewTransport,
		"grpc": transportgrpc.NewTransport,
	}
)
