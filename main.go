package main

import (
	"fmt"
	"log"

	"go.etcd.io/etcd/client"
	"github.com/clicktherapeutics/ct-dns/pkg/etcd"
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

	etcdCli := etcd.NewAPIClient(client.NewKeysAPI(c))

	// Set key "/foo" to value "bar".
	err = etcdCli.CreateOrSet("/foo", "bar2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("successfully create/set key %q with value %q", "/foo", "bar")

	val, err := etcdCli.Get("/foo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("key has %q value\n", val)
}
