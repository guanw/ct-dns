package grpc

import (
	"context"

	pb "github.com/guanw/ct-dns/pkg/grpc/proto-gen"
	"github.com/guanw/ct-dns/pkg/store"
)

// DNSServer implements pb.DnsServer
type DNSServer struct {
	Store   store.Store
	Metrics *Metrics
}

// NewServer creates new DnsServer
func NewServer(store store.Store, metrics *Metrics) pb.DnsServer {
	return &DNSServer{
		Store:   store,
		Metrics: metrics,
	}
}

// GetService implements DnsServer.GetService
func (s *DNSServer) GetService(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	serviceName := req.GetServiceName()
	hosts, err := s.Store.GetService(serviceName)
	if err != nil {
		s.Metrics.GetServiceFailure.Inc()
	} else {
		s.Metrics.GetServiceSuccess.Inc()
	}
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
	if err != nil {
		s.Metrics.PostServiceFailure.Inc()
	} else {
		s.Metrics.PostServiceSuccess.Inc()
	}
	return &pb.PostResponse{}, err
}
