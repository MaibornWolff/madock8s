##### Build Sample metrics 

```
cd $GOPATH/src/github.com/MaibornWolff/maDocK8s/exporter/prometheus/sample-metrics
docker build -t maibornwolff/madock8s_sample_metrics .
```

##### Deploy SampleMetrics

```
kubectl apply -f sample-metrics.yaml
```

Check Prometheus exporter logs, it must be notified by the controller.
