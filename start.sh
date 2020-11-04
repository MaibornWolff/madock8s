#! /bin/bash

set -e
echo "Hello madock8s Developer!"

# minikube start --kubernetes-version=1.15.9 --driver=virtualbox &&
eval $(minikube docker-env)


echo "Building docker images"
docker-compose build --parallel

echo "Upgrading charts"

mdstorageHelmPath=core/services/mdstorage/helm
helm upgrade -i mdstorage -f $mdstorageHelmPath/values.yaml $mdstorageHelmPath

dashboardHelmPath=dashboard
helm upgrade -i madock8s-dashboard -f $dashboardHelmPath/values.yaml $dashboardHelmPath

controllerHelmPath=core/controller/helm
helm upgrade -i madock8s-controller -f $controllerHelmPath/values.yaml $controllerHelmPath

envExporterHelmPath=exporter/env/helm
helm upgrade -i env-exporter -f $envExporterHelmPath/values.yaml $envExporterHelmPath

githubExporterHelmPath=exporter/github/helm
helm upgrade -i github-exporter -f $githubExporterHelmPath/values.yaml $githubExporterHelmPath

gitlabExporterHelmPath=exporter/gitlab/helm
helm upgrade -i gitlab-exporter -f $gitlabExporterHelmPath/values.yaml $gitlabExporterHelmPath

prometheusExporterHelmPath=exporter/prometheus/helm
helm upgrade -i prometheus-exporter -f $prometheusExporterHelmPath/values.yaml $prometheusExporterHelmPath

swaggerExporterHelmPath=exporter/swagger/helm
helm upgrade -i swagger-exporter -f $swaggerExporterHelmPath/values.yaml $swaggerExporterHelmPath

versionExporterHelmPath=exporter/version/helm
helm upgrade -i version-exporter -f $versionExporterHelmPath/values.yaml $versionExporterHelmPath

echo "Applying sample metrics"
sampleMetricsPath=sample-metrics/yaml
kubectl apply -f $sampleMetricsPath/configmap.yaml
kubectl apply -f $sampleMetricsPath/secret.yaml
kubectl apply -f $sampleMetricsPath/service.yaml
kubectl apply -f $sampleMetricsPath/deployment.yaml
