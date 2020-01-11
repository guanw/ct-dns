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
install:
	env GO111MODULE=on go build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)
test:
	env GO111MODULE=on $(GOTEST) -v ./...
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

docker-build:
	docker build -t ct-dns .
docker-run:
	docker run -d --rm -p 8080:8080 -p 50051:50051 ct-dns