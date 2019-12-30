package grpc

import (
	"context"

	pb "github.com/clicktherapeutics/ct-dns/pkg/grpc/proto-gen"
	"github.com/clicktherapeutics/ct-dns/pkg/store"
)

// DNSServer implements pb.DnsServer
type DNSServer struct {
	Store store.Store
}

// NewServer creates new DnsServer
func NewServer(store store.Store) pb.DnsServer {
	return &DNSServer{
		Store: store,
	}
}

// GetService implements DnsServer.GetService
func (s *DNSServer) GetService(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	serviceName := req.GetServiceName()
	hosts, err := s.Store.GetService(serviceName)
	return &pb.GetResponse{
		Hosts: hosts,
	}, err
}

// PostService implements DnsServer.PostService
func (s *DNSServer) PostService(ctx context.Context, req *pb.PostRequest) (*pb.PostResponse, error) {
	err := s.Store.UpdateService(
		req.GetServiceName(),
		req.GetOperation(),
		req.GetHost(),
	)
	return &pb.PostResponse{}, err
}
