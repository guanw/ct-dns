package grpc

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	pb "github.com/clicktherapeutics/ct-dns/pkg/grpc/proto-gen"
	"github.com/clicktherapeutics/ct-dns/pkg/store"
	"github.com/clicktherapeutics/ct-dns/pkg/store/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func initialize(store store.Store) {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterDnsServer(s, NewServer(store))
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(string, time.Duration) (net.Conn, error) {
	return lis.Dial()
}

func Test_GetServiceSucceed(t *testing.T) {
	store := &mocks.Store{}
	store.On("GetService", "valid-service").Return([]string{"192.0.0.1"}, nil)
	initialize(store)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
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
}

func Test_PostService(t *testing.T) {
	store := &mocks.Store{}
	store.On("UpdateService", "valid-service", "add", "192.0.0.1").Return(nil)
	initialize(store)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
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
}
