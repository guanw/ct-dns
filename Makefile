# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=ct_dns_binary
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
lint:
	revive -formatter friendly -exclude github.com/guanw/ct-dns/vendor/... github.com/guanw/ct-dns/...
fmt:
	go fmt ./...
	golangci-lint run ./...
install:
	@which dep > /dev/null || curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	dep ensure -vendor-only

install-cobra:
	go install github.com/guanw/ct-dns

build:
	$(GOBUILD) -o $(BINARY_NAME) -v
# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)
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
	echo "pass the public ip above to --etcd-endpoints"

dynamodb-single-cluster:
	docker stop dns-dynamodb && docker rm dns-dynamodb && docker run -d -it -p 8000:8000 --name dns-dynamodb dwmkerr/dynamodb -sharedDb

redis-single-cluster:
	docker stop dns-redis && docker rm dns-redis && docker run -d -p 6379:6379 --name dns-redis redis
	docker exec -it dns-redis sh

docker-build:
	docker build -t ct-dns .

docker-run:
	docker run -d --rm -p 8080:8080 -p 50051:50051 ct-dns