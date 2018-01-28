package micro

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-log/log"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/transport"
)

func NewService(sc *ServiceConfig, opts ...Option) Service {
	s := &service{
		sc:   sc,
		opts: newOptions(sc, opts...),
	}
	return s
}

func newOptions(sc *ServiceConfig, opts ...Option) Options {
	b := newBroker(sc)

	t := newTransport(sc)

	r := newRegistry(sc)

	sel := newSelector(sc, selector.Registry(r))

	c := newClient(
		sc,
		client.Broker(b),
		client.Transport(t),
		client.Selector(sel),
		client.Registry(r),
	)
	c = &clientWrapper{
		c,
		metadata.Metadata{
			HeaderPrefix + "From-Service": sc.ServerName,
		},
	}

	s := newServer(
		sc,
		server.Broker(b),
		server.Transport(t),
		server.Registry(r),
	)

	opt := Options{
		Broker:    b,
		Client:    c,
		Server:    s,
		Transport: t,
		Registry:  r,
		Context:   context.Background(),
	}

	if sc.RegisterInterval != "" {
		d, err := time.ParseDuration(sc.RegisterInterval)
		if err != nil {
			panic(fmt.Sprintf("failed to parse RegisterInterval: %s", sc.RegisterInterval))
		}
		opts = append(opts, RegisterInterval(d))
	}

	if sc.RegisterTTL != "" {
		d, err := time.ParseDuration(sc.RegisterTTL)
		if err != nil {
			panic(fmt.Sprintf("failed to parse RegisterTTL: %s", sc.RegisterTTL))
		}
		opts = append(opts, RegisterTTL(d))
	}

	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func newBroker(sc *ServiceConfig) broker.Broker {
	var opts []broker.Option
	if len(sc.BrokerAddress) > 0 {
		opts = append(opts, broker.Addrs(sc.BrokerAddress...))
	}
	if sc.Broker != "" {
		b := DefaultBrokers[sc.Broker]
		if b == nil {
			panic(fmt.Sprintf("unknwon broker %s", sc.Broker))
		}
		return b(opts...)
	}
	return DefaultBrokers[defaultBroker](opts...)
}

func newTransport(sc *ServiceConfig) transport.Transport {
	var opts []transport.Option
	if len(sc.TransportAddress) > 0 {
		opts = append(opts, transport.Addrs(sc.TransportAddress...))
	}
	if sc.Transport != "" {
		t := DefaultTransports[sc.Transport]
		if t == nil {
			panic(fmt.Sprintf("unknwon transport %s", sc.Transport))
		}
		return t(opts...)
	}
	return DefaultTransports[defaultTransport](opts...)
}

func newSelector(sc *ServiceConfig, opts ...selector.Option) selector.Selector {
	if sc.Selector != "" {
		s := DefaultSelectors[sc.Selector]
		if s == nil {
			panic(fmt.Sprintf("unknwon selector %s", sc.Selector))
		}
		return s(opts...)
	}
	return DefaultSelectors[defaultSelector](opts...)
}

func newRegistry(sc *ServiceConfig) registry.Registry {
	var opts []registry.Option
	if len(sc.RegistryAddress) > 0 {
		opts = append(opts, registry.Addrs(sc.RegistryAddress...))
	}
	if sc.Registry != "" {
		r := DefaultRegistries[sc.Registry]
		if r == nil {
			panic(fmt.Sprintf("unknwon registry %s", sc.Registry))
		}
		return r(opts...)
	}
	return DefaultRegistries[defaultRegistry](opts...)
}

func newClient(sc *ServiceConfig, opts ...client.Option) client.Client {
	if sc.ClientPoolSize > 0 {
		opts = append(opts, client.PoolSize(sc.ClientPoolSize))
	}
	if sc.ClientPoolTTL != "" {
		d, err := time.ParseDuration(sc.ClientPoolTTL)
		if err != nil {
			panic(fmt.Sprintf("failed to parse client_pool_ttl: %s", sc.ClientPoolTTL))
		}
		opts = append(opts, client.PoolTTL(d))
	}

	if sc.ClientRequestTimeout != "" {
		d, err := time.ParseDuration(sc.ClientRequestTimeout)
		if err != nil {
			panic(fmt.Sprintf("failed to parse client_request_timeout: %s", sc.ClientRequestTimeout))
		}
		opts = append(opts, client.RequestTimeout(d))
	}

	if sc.ClientRetries > 0 {
		opts = append(opts, client.Retries(sc.ClientRetries))
	}

	if sc.Client != "" {
		c := DefaultClients[sc.Client]
		if c == nil {
			panic(fmt.Sprintf("unknwon client %s", sc.Client))
		}
		return c(opts...)
	}
	return DefaultClients[defaultClient](opts...)
}

func newServer(sc *ServiceConfig, opts ...server.Option) server.Server {
	if sc.ServerName != "" {
		opts = append(opts, server.Name(sc.ServerName))
	}

	if sc.ServerVersion != "" {
		opts = append(opts, server.Version(sc.ServerVersion))
	}

	if sc.ServerAddress != "" {
		opts = append(opts, server.Address(sc.ServerAddress))
	}

	if sc.ServerAdvertise != "" {
		opts = append(opts, server.Advertise(sc.ServerAdvertise))
	}

	if sc.ServerID != "" {
		opts = append(opts, server.Id(sc.ServerID))
	}

	if len(sc.ServerMetaData) > 0 {
		metadata := make(map[string]string)
		for _, d := range sc.ServerMetaData {
			var key, val string
			parts := strings.Split(d, "=")
			key = parts[0]
			if len(parts) > 1 {
				val = strings.Join(parts[1:], "=")
			}
			metadata[key] = val
		}
		opts = append(opts, server.Metadata(metadata))
	}

	if sc.Server != "" {
		s := DefaultServers[sc.Server]
		if s == nil {
			panic(fmt.Sprintf("unknwon server %s", sc.Server))
		}
		return s(opts...)
	}
	return DefaultServers[defaultServer](opts...)
}

type service struct {
	sc   *ServiceConfig
	opts Options
}

func (s *service) run(exit chan bool) {
	if s.opts.RegisterInterval <= time.Duration(0) {
		return
	}

	t := time.NewTicker(s.opts.RegisterInterval)

	for {
		select {
		case <-t.C:
			err := s.opts.Server.Register()
			if err != nil {
				log.Log("service run Server.Register error: ", err)
			}
		case <-exit:
			t.Stop()
			return
		}
	}
}

func (s *service) Init(opts ...Option) {}

func (s *service) Options() Options {
	return s.opts
}

func (s *service) Client() client.Client {
	return s.opts.Client
}

func (s *service) Server() server.Server {
	return s.opts.Server
}

func (s *service) String() string {
	return "go-micro"
}

func (s *service) Start() error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	if err := s.opts.Server.Start(); err != nil {
		return err
	}

	if err := s.opts.Server.Register(); err != nil {
		return err
	}

	// for _, fn := range s.opts.AfterStart {
	// 	if err := fn(); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (s *service) Stop() error {
	var gerr error

	// for _, fn := range s.opts.BeforeStop {
	// 	if err := fn(); err != nil {
	// 		gerr = err
	// 	}
	// }

	if err := s.opts.Server.Deregister(); err != nil {
		return err
	}

	if err := s.opts.Server.Stop(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	return gerr
}

func (s *service) Run() error {
	if err := s.Start(); err != nil {
		return err
	}

	// start reg loop
	ex := make(chan bool)
	go s.run(ex)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	select {
	// wait on kill signal
	case <-ch:
	// wait on context cancel
	case <-s.opts.Context.Done():
	}

	// exit reg loop
	close(ex)

	if err := s.Stop(); err != nil {
		return err
	}

	return nil
}
