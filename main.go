package main

import (
	"fmt"
	"log"

	"github.com/clicktherapeutics/ct-dns/pkg/etcd"
	dHttp "github.com/clicktherapeutics/ct-dns/pkg/http"
	"github.com/gorilla/mux"
	"go.etcd.io/etcd/client"
	"net/http"
)

func main() {
	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:5001"},
	}

	c, err := client.New(cfg)
	if err != nil {
		// handle error
		fmt.Println("Cannot initialize the etcd client")
	}

	etcdCli := etcd.NewClient(client.NewKeysAPI(c))

	r := mux.NewRouter()
	httpHandler := dHttp.NewHandler(etcdCli)
	httpHandler.RegisterRoutes(r)
	if err = http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatal(err)
	}
}
