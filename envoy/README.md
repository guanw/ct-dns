# Run ct-dns with envoy locally

### install envoy

```
$ brew tap tetratelabs/getenvoy
$ brew install envoy
```

### start local redis cluster

```
$ make redis-single-cluster
```

### start ct-dns

```
$ go run main.go --storage-type=redis
```

### start an sample application on 8081

```
$ cd upstream/

$ virtualenv env --python=python2.7
$ source env/bin/activate
$ pip install -r requirements.txt

$ python server.py -p 8081
```

### manually register sample application by POST to ct-dns

```
POST http://localhost:8080/api/service HTTP/1.1
Content-Type: application/json

{
    "serviceName": "dummy-service",
    "operation": "add",
    "host": "0.0.0.0:8081"
}
```

### once it's registered, envoy EDS cluster will start health checking on this host&port

If you do

```
POST http://localhost:8080/api/service HTTP/1.1
Content-Type: application/json

{
    "serviceName": "dummy-service",
    "operation": "delete",
    "host": "0.0.0.0:8081"
}
```

You should see health checking stopped
