# MaDocK8s - Markdown Kubernetes Framework

## Prerequisites

Installed tools:
- minikube
- helm
- docker with docker-compose


## HowTo:

1. Start minikube:
```
    minikube start --kubernetes-version=1.15.9 --driver=virtualbox
```

- Mandatory for the **first** run. Execute:
```
    helm init 
    chmod u+x ./start.sh ./clean.sh
```

2. Execute `./start.sh`. The script will build all required docker images and install helm charts.

3. Check the dashboard:
    - Print url of dashboard to console 
```
    minikube service --url madock8s-dashboard
```

    - Open the url in your favourite browser;
    - Press "View Documentation".

4. To reset the cluster, execute `./reset.sh`

## Sample App

There is an example app in `sample-metrics` directory. It has configurations for **all** MaDocK8s exporters. Check `metadata.annotations` in  `sample-metrics/yaml/deployment.yaml`. 

If you edit the deployment and apply the changes, the updated documentation in the dashboard will be available after a delay of approx. 5 seconds. 
