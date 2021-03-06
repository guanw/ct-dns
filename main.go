package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"net/http"

	"github.com/gorilla/mux"
	config "github.com/guanw/ct-dns/cmd"
	dns "github.com/guanw/ct-dns/pkg/grpc"
	pb "github.com/guanw/ct-dns/pkg/grpc/proto-gen"
	ctHttp "github.com/guanw/ct-dns/pkg/http"
	"github.com/guanw/ct-dns/pkg/logging"
	ctStore "github.com/guanw/ct-dns/pkg/store"
	"github.com/guanw/ct-dns/plugins/storage"
	"github.com/guanw/ct-dns/plugins/storage/dynamodb"
	"github.com/guanw/ct-dns/plugins/storage/etcd"
	"github.com/guanw/ct-dns/plugins/storage/redis"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
			f := storage.NewFactory(v, cfg)
			client, err := f.Initialize()
			if err != nil {
				return errors.Wrap(err, "Failed to start storage client")
			}
			store := ctStore.NewStore(client)
			// TODO move 5 to config/from flag
			retryStore := ctStore.NewRetryHandler(5, store, ctStore.InitializeMetrics())
			dnsServer := dns.NewServer(retryStore, dns.InitializeMetrics())
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
			logging.GetLogger().Printf("grpc server listening at port %s", cfg.GRPCPort)

			r := mux.NewRouter()
			httpHandler := ctHttp.NewHandler(retryStore, ctHttp.InitializeMetrics())
			httpHandler.RegisterRoutes(r)

			r.Handle("/metrics", promhttp.Handler())
			logging.GetLogger().Printf("http server listening at port %s", cfg.HTTPPort)
			return http.ListenAndServe("0.0.0.0:"+cfg.HTTPPort, r)
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
	etcd.AddFlags(flagSet)
	redis.AddFlags(flagSet)
	storage.AddFlags(flagSet)

	command.Flags().AddGoFlagSet(flagSet)
	v.BindPFlags(command.Flags())
}
