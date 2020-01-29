# ct-dns - Service Discovery System

<img src="https://scionplu.sirv.com/ct-dns.jpg" width="614" height="264" alt="" />

ct-dns aims to provide easy deployment of service discovery service.

1. It supports following protocols

- http
- grpc

2. It supports following storage options

- dynamodb
- etcd
- redis
- memory (mainly for testing it out)

3. It supports integrating with envoy as eds cluster with examples.

# Development

`$make install`

`$make test`

`$make run`

You will see both http and grpc server up and running like following in console:

<img src="https://scionplu.sirv.com/Images/server.png" width="300" height="35" alt="" />

# Start up local etcd cluster:

1. single node: `$make etcd-single-node`

2. kuberneters three-node cluster: `$make etcd-kube`

# Start up local dynamodb cluster:

`$make dynamodb-single-cluster`
