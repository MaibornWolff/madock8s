### Howto deploy to local minikube

requirements: docker, minikube, helm 2

```
# start minikube
minikube start --kubernetes-version=1.15.9
eval $(minikube docker-env)

# build controller
cd $GOPATH/src/github.com/MaibornWolff/maDocK8s/core/controller
docker build -t maibornwolff/madock8s-controller .

# deploy controller
cd $GOPATH/src/github.com/MaibornWolff/maDocK8s/core/controller/helm/charts/madock8s-controller
helm init # (for the first time only)
helm del --purge madock8s-controller # (to remove older version)
helm install -n madock8s-controller -f values.yaml .

# build an install dummy exporter
cd $GOPATH/src/github.com/MaibornWolff/maDocK8s/grpc/server/
docker build -t maibornwolff/grpctest .
k apply -f server.yaml

# watch controller
kubectl logs $(k get pods | grep madock8s-controller | grep -Eo '^[^ ]+') -f

# watch dummy exporter
kubectl logs $(k get pods | grep grpcserver | grep -Eo '^[^ ]+') -f

# deploy test service 
cd $GOPATH/src/github.com/MaibornWolff/maDocK8s/core/controller/helm/charts/madock8s-controller
kubectl apply -f nginx.yaml
# The depoyment is detected by the controller and the dummy exporter are notified (see logs)

# If you change the deployment now (e.g. increase replicas to 2) and deploy again, controller and dummy exporter are notifiedn again
```