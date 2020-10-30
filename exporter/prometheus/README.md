# Prometheus Exporter

Prometheus Exporter is the service that will fetch and read metrics from given microservices.


## Configuration

To retrieve the list of available metrics, Prometheus Exporter uses service's clusterIP:80/metrics endpoint.
TODO: Move endpoint configuration to ENV vars

### Target Deployment

To export metrics, the target deployment must provide the following annotation in metadata:

```
madock8s.exporter/prometheusExporter: true
```

The value is a string-boolean with semantic meaning "exporter is enabled".
For now, the value can be any non-empty string. 


### Delete Mode

Configure behaviour for md-files when target deployment is deleted.
Simply update helm/values.yaml/env.deletionStrategy to one of the following values:
- IGNORE - keep existing md-files without modification;
- UPDATE - add "Deployments of the service were deleted on timestamp" to existing md-file (default);
- DELETE - remove the respective md-file from storage.
