# start local kubernetes cluster

```
$ minikube start
```

# enable minikube ingress addon

```
$ minikube addons enable ingress
```

# apply deployment.yml

```
$ kubectl apply -f deployment.yml
```

# verify the pod/svc is running

```
$ kubectl get pods
```

You should see something similar to the following

<img src="https://scionplu.sirv.com/kube.png" width="300" height="70" alt="" />

# Port forward so you can test curl using localhost

```
$ kubectl port-forward <pod-name> 8080:8080 50050:50050
```

# or directly by

```
curl <minikube-ip>/api/health
```
