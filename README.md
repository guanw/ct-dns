# ct-dns - Distributed Service Discovery System

ct-dns aims to provide easy deployment of service discovery service. It exposes

- http server
- grpc server

for interfacing and

- dynamodb
- etcd
- memory

for storage options

# Development

`$make install`

`$make test`

`$make run`

You will see both http and grpc server up and running like following in console:

<img src="https://scionplu.sirv.com/Images/server.png" width="500" height="30" alt="" />

# Start up local etcd cluster:

1. single node: `$make etcd-single-node`

2. kuberneters three-node cluster: `$make etcd-kube`

# Start up local dynamodb cluster:

`$make dynamodb-single-cluster`
