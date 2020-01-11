package main

import (
	"fmt"
	"log"
	"net"

	"net/http"

	"github.com/gorilla/mux"
	config "github.com/guanw/ct-dns/cmd"
	dns "github.com/guanw/ct-dns/pkg/grpc"
	pb "github.com/guanw/ct-dns/pkg/grpc/proto-gen"
	ctHttp "github.com/guanw/ct-dns/pkg/http"
	"github.com/guanw/ct-dns/pkg/store"
	"github.com/guanw/ct-dns/plugins/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	cfg := config.ReadConfig("./config/")

	f := storage.NewFactory()
	client, err := f.Initialize("redis")
	if err != nil {
		log.Fatalf("Failed to start storage client: %v", err)
	}
	store := store.NewStore(client)

	dnsServer := dns.NewServer(store)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	healthServer := health.NewServer()
	healthServer.SetServingStatus("ct-dns", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
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
