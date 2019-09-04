package main

import (
	"log"

	"github.com/guoxf/micro/go-micro"
)

func main() {
	log.SetFlags(log.Llongfile | log.Ltime)
	sc := micro.NewServiceConfig("./service.yaml")
	log.Printf("%#v\n", *sc)
	s := micro.NewService(sc)
	log.Println(s.Client().Options().Registry.String())
	if err := s.Run(); err != nil {
		panic(err)
	}
}
