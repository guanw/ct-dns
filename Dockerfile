FROM golang:1.9

WORKDIR /go/src/github.com/guanw/ct-dns
COPY . .
RUN make build-linux
EXPOSE 8080
EXPOSE 50051
CMD ["/go/src/github.com/guanw/ct-dns/ct_dns_binary_unix"]