package main

import (
	"fmt"
	"log"

	"net/http"

	config "github.com/clicktherapeutics/ct-dns/cmd"
	"github.com/clicktherapeutics/ct-dns/pkg/etcd"
	ctHttp "github.com/clicktherapeutics/ct-dns/pkg/http"
	"github.com/gorilla/mux"
	"go.etcd.io/etcd/client"
)

func main() {
	cfg := config.ReadConfig("./config/")
	var endpoints []string
	for _, v := range cfg.Etcd {
		endpoints = append(endpoints, v.Host+":"+v.Port)
	}
	etcdCfg := client.Config{
		Endpoints: endpoints,
	}

	c, err := client.New(etcdCfg)
	if err != nil {
		// handle error
		fmt.Println("Cannot initialize the etcd client")
	}

	etcdCli := etcd.NewClient(client.NewKeysAPI(c))

	r := mux.NewRouter()
	httpHandler := ctHttp.NewHandler(etcdCli)
	httpHandler.RegisterRoutes(r)
	if err = http.ListenAndServe("0.0.0.0:"+cfg.HTTPPort, r); err != nil {
		log.Fatal(err)
	}
}
