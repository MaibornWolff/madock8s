#! /bin/bash

echo "Purging helm charts"

arg=""
helm_version=$(helm version --short)
if [[ $helm_version == *"v2."* ]]; then
  arg=--purge
fi

helm del $arg madock8s-controller
helm del $arg prometheus-exporter
helm del $arg gitlab-exporter
helm del $arg github-exporter
helm del $arg env-exporter
helm del $arg swagger-exporter
helm del $arg version-exporter
helm del $arg mdstorage
helm del $arg madock8s-dashboard

echo "Deleting sample-metrics"
sampleMetricsPath=sample-metrics/yaml
kubectl delete -f $sampleMetricsPath/configmap.yaml
kubectl delete -f $sampleMetricsPath/secret.yaml
kubectl delete -f $sampleMetricsPath/service.yaml
kubectl delete -f $sampleMetricsPath/deployment.yaml

echo "Deleting daux-generator pods"
kubectl delete pods --selector=job-name=daux-generator
