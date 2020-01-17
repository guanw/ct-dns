package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"net/http"

	"github.com/gorilla/mux"
	config "github.com/guanw/ct-dns/cmd"
	dns "github.com/guanw/ct-dns/pkg/grpc"
	pb "github.com/guanw/ct-dns/pkg/grpc/proto-gen"
	ctHttp "github.com/guanw/ct-dns/pkg/http"
	"github.com/guanw/ct-dns/pkg/store"
	"github.com/guanw/ct-dns/plugins/storage"
	"github.com/guanw/ct-dns/plugins/storage/dynamodb"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	cfg := config.ReadConfig("./config/")
	v := viper.New()
	command := &cobra.Command{
		Use:   "ct-dns",
		Short: "ct-dns register and update host information for specific service",
		Long:  `ct-dns register and update host information for specific service, User can configure different storage types using terminal flag`,
		RunE: func(cmd *cobra.Command, args []string) error {
			f := storage.NewFactory()
			client, err := f.Initialize(v, "dynamodb")
			if err != nil {
				return errors.Wrap(err, "Failed to start storage client")
			}
			store := store.NewStore(client)

			dnsServer := dns.NewServer(store)
			lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
			if err != nil {
				return errors.Wrap(err, "Failed to listen")
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
				return err
			}
			return nil
		},
	}
	AddFlags(v, command)
	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// AddFlags applies binding flags to initialize app
func AddFlags(v *viper.Viper, command *cobra.Command) {
	flagSet := new(flag.FlagSet)

	dynamodb.AddFlags(flagSet)

	command.Flags().AddGoFlagSet(flagSet)
	v.BindPFlags(command.Flags())
}
