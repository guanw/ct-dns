package grpc

import (
	"context"

	pb "github.com/clicktherapeutics/ct-dns/pkg/grpc/proto-gen"
)

// DNSServer implements pb.DnsServer
type DNSServer struct{}

// NewServer creates new DnsServer
func NewServer() pb.DnsServer {
	return &DNSServer{}
}

// GetService implements DnsServer.GetService
func (s *DNSServer) GetService(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	return &pb.GetResponse{}, nil
}

// PostService implements DnsServer.PostService
func (s *DNSServer) PostService(ctx context.Context, req *pb.PostRequest) (*pb.PostResponse, error) {
	return &pb.PostResponse{}, nil
}
