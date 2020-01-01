package main

import (
	"fmt"
	"log"
	"net"

	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	config "github.com/clicktherapeutics/ct-dns/cmd"
	dns "github.com/clicktherapeutics/ct-dns/pkg/grpc"
	pb "github.com/clicktherapeutics/ct-dns/pkg/grpc/proto-gen"
	ctHttp "github.com/clicktherapeutics/ct-dns/pkg/http"
	ctdynamodb "github.com/clicktherapeutics/ct-dns/pkg/storage/dynamodb"
	"github.com/clicktherapeutics/ct-dns/pkg/store"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	cfg := config.ReadConfig("./config/")
	var endpoints []string
	for _, v := range cfg.Etcd {
		endpoints = append(endpoints, "http://"+v.Host+":"+v.Port)
	}
	// etcdCfg := client.Config{
	// 	Endpoints: endpoints,
	// }

	// c, err := client.New(etcdCfg)
	// if err != nil {
	// 	// handle error
	// 	log.Fatalf(errors.Wrap(err, "Cannot initialize the etcd client").Error())
	// }

	// etcdCli := etcd.NewClient(client.NewKeysAPI(c))
	// store := store.NewStore(etcdCli)

	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:8000"),
	}))
	db := dynamodb.New(sess)

	store := store.NewStore(ctdynamodb.NewClient(db))

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
