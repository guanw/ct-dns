package main

import (
	"fmt"
	"log"
	"net"

	"net/http"

	config "github.com/clicktherapeutics/ct-dns/cmd"
	"github.com/clicktherapeutics/ct-dns/pkg/etcd"
	dns "github.com/clicktherapeutics/ct-dns/pkg/grpc"
	pb "github.com/clicktherapeutics/ct-dns/pkg/grpc/proto-gen"
	ctHttp "github.com/clicktherapeutics/ct-dns/pkg/http"
	"github.com/clicktherapeutics/ct-dns/pkg/store"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/client"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.ReadConfig("./config/")
	var endpoints []string
	for _, v := range cfg.Etcd {
		endpoints = append(endpoints, "http://"+v.Host+":"+v.Port)
	}
	etcdCfg := client.Config{
		Endpoints: endpoints,
	}

	c, err := client.New(etcdCfg)
	if err != nil {
		// handle error
		log.Fatalf(errors.Wrap(err, "Cannot initialize the etcd client").Error())
	}

	etcdCli := etcd.NewClient(client.NewKeysAPI(c))
	store := store.NewStore(etcdCli)

	dnsServer := dns.NewServer(store)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterDnsServer(grpcServer, dnsServer)
	go grpcServer.Serve(lis)
	defer grpcServer.Stop()
	log.Printf("grpc server listening at port %s", cfg.GRPCPort)

	r := mux.NewRouter()
	httpHandler := ctHttp.NewHandler(store)
	httpHandler.RegisterRoutes(r)
	log.Printf("http server listening at port %s", cfg.HTTPPort)
	if err = http.ListenAndServe("0.0.0.0:"+cfg.HTTPPort, r); err != nil {
		log.Fatal(err)
	}
}
