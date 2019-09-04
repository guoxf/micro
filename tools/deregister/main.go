package main

import (
	"flag"
	"fmt"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-plugins/registry/etcd"
)

func main() {
	registryName := flag.String("r", "consul", "registry name")
	host := flag.String("h", "127.0.0.1:8500", "registry address")
	serviceName := flag.String("s", "", "service name")
	flag.Parse()
	fmt.Println(*registryName, *host, *serviceName)
	r := newRegistry(*registryName, *host)
	fmt.Println(r.String())
	if *serviceName != "" {
		deregisterByServiceName(r, *serviceName)
		return
	}

	for {
		services, err := selectService(r)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for i := range services {
			deregisterNode(r, services[i])
		}

	}
}

func deregisterNode(r registry.Registry, s *registry.Service) {
	for i, node := range s.Nodes {
		fmt.Printf("%d. %s %s\n", i, node.Address, node.Id)
	}
	fmt.Println("please select node")
	var nodeIndex int
	_, err := fmt.Scanln(&nodeIndex)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = r.Deregister(&registry.Service{
		Name:    s.Name,
		Version: s.Version,
		Nodes:   []*registry.Node{s.Nodes[nodeIndex]},
	})
	if err != nil {
		fmt.Println(err)
	}
}

func deregisterService(r registry.Registry, s *registry.Service) {
	for i, node := range s.Nodes {
		fmt.Printf("%d. deregister %s %s\n", i, node.Address, node.Id)
		err := r.Deregister(&registry.Service{
			Name:    s.Name,
			Version: s.Version,
			Nodes:   []*registry.Node{node},
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}

func deregisterByServiceName(r registry.Registry, name string) {
	services, err := r.GetService(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(services)
	for i := range services {
		deregisterService(r, services[i])
	}
}

func selectService(r registry.Registry) ([]*registry.Service, error) {
	services, err := r.ListServices()
	if err != nil {
		panic(err)
	}

	for i := range services {
		fmt.Printf("%d. %s\n", i, services[i].Name)
	}
	fmt.Println("please select service")
	var serviceIndex int
	_, err = fmt.Scanln(&serviceIndex)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return r.GetService(services[serviceIndex].Name)
}

func newRegistry(name, host string) registry.Registry {
	opt := registry.Addrs(host)
	switch name {
	case "etcd":
		return etcd.NewRegistry(opt)
	default:
		return consul.NewRegistry(opt)
	}
}
