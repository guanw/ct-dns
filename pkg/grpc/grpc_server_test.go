package grpc

import (
	"context"
	"errors"
	"net"
	"testing"

	pb "github.com/guanw/ct-dns/pkg/grpc/proto-gen"
	"github.com/guanw/ct-dns/pkg/logging"
	"github.com/guanw/ct-dns/pkg/store"
	"github.com/guanw/ct-dns/pkg/store/mocks"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var (
	lis     *bufconn.Listener
	metrics = InitializeMetrics()
)

func initialize(store store.Store) {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterDnsServer(s, NewServer(store, metrics))
	go func() {
		if err := s.Serve(lis); err != nil {
			logging.GetLogger().Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func Test_GetServiceSucceed(t *testing.T) {
	store := &mocks.Store{}
	store.On("GetService", "valid-service").Return([]string{"192.0.0.1"}, nil)
	initialize(store)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewDnsClient(conn)
	resp, err := client.GetService(ctx, &pb.GetRequest{
		ServiceName: "valid-service",
	})
	assert.NoError(t, err)
	assert.Equal(t, []string{"192.0.0.1"}, resp.GetHosts())
	//TODO replace with testutil.CollectAndCount with new release
	assert.Equal(t, 1.0, testutil.ToFloat64(metrics.GetServiceSuccess))
}

func Test_GetServiceFail(t *testing.T) {
	store := &mocks.Store{}
	store.On("GetService", "error-service").Return(nil, errors.New("get service failed"))
	initialize(store)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewDnsClient(conn)
	_, err = client.GetService(ctx, &pb.GetRequest{
		ServiceName: "error-service",
	})
	assert.Error(t, err)
	assert.Equal(t, 1.0, testutil.ToFloat64(metrics.GetServiceFailure))
}

func Test_PostServiceSucceed(t *testing.T) {
	store := &mocks.Store{}
	store.On("UpdateService", "valid-service", "add", "192.0.0.1").Return(nil)
	initialize(store)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewDnsClient(conn)
	_, err = client.PostService(ctx, &pb.PostRequest{
		ServiceName: "valid-service",
		Operation:   "add",
		Host:        "192.0.0.1",
	})
	assert.NoError(t, err)
	assert.Equal(t, 1.0, testutil.ToFloat64(metrics.PostServiceSuccess))
}

func Test_PostServiceFail(t *testing.T) {
	store := &mocks.Store{}
	store.On("UpdateService", "error-service", "add", "192.0.0.1").Return(errors.New("service update failed"))
	initialize(store)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewDnsClient(conn)
	_, err = client.PostService(ctx, &pb.PostRequest{
		ServiceName: "error-service",
		Operation:   "add",
		Host:        "192.0.0.1",
	})
	assert.Error(t, err)
	assert.Equal(t, 1.0, testutil.ToFloat64(metrics.PostServiceFailure))
}
