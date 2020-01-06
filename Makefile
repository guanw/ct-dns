# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=ct_dns_binary
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
fmt:
	go fmt ./...
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
protoc:
	cd IDL/proto && protoc -I . dns.proto --go_out=plugins=grpc:.

etcd-single-node:
	cd etcd && chmod +x etcd.sh && ./etcd.sh

etcd-kube:
	cd etcd && kubectl apply -f etcd-sts.yaml
	minikube tunnel &
	kubectl get all -n default | grep etcd-client
	echo "replace config/development.yaml host with the public ip above"

dynamodb-single-cluster:
	docker run -d -it -p 8000:8000 dwmkerr/dynamodb -sharedDb
# run:
# 	$(GOBUILD) -o $(BINARY_NAME) -v ./...
# 	./$(BINARY_NAME)


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/clicktherapeutics/ct-dns golang:latest go build -o "$(BINARY_UNIX)" -v