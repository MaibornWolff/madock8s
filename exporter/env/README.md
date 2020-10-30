# Environment Exporter

Environment Exporter collects environment variables declared for all containers in Kubernetes Deployments, available via Kubernetes API.


Supported Declarations of Environment Variables:

- direct declaration in configuration file
- fieldRef in configuration file
- resourceFieldRef in configuration file
- referenced environment variables from configmaps
- referenced environment variables that from secrets


## Configuration

### Target Deployment

To export environment variables, the target deployment must provide the following annotation in metadata:

```
madock8s.exporter/envExporter: sample-metrics
```

The value is the name of deployment of the target service.


### Delete Mode

Configure behaviour for md-files when target deployment is deleted.
Simply update helm/values.yaml/env.deletionStrategy to one of the following values:
- IGNORE - keep existing md-files without modification;
- UPDATE - add "Deployments of the service were deleted on timestamp" to existing md-file (default);
- DELETE - remove the respective md-file from storage.
