# Version Exporter

Version Exporter retrieves details about images used in containers of the service.

## Configuration

### Service configuration
Target deployment has to provide the following annotation in metadata:
```
    madock8s.exporter/versionExporter: true
```
The value semantically suggests that Version Exporter is should be boolean, indicating that version exporter is enabled for this deployment. However, any value (also `false`) is suffiecient to enable version exporter for the deployment.

Whenever there is an update to deployment that has specified annotation, MaDocK8s-Controller will notify Version Exporter. 

### Workmodes
There are 2 workmodes of Version Exporter. This can be configured in helm/values.yaml by setting boolean `env.labeledOnly`.

Allowed values:
- __false__ (default): collect information about images for containers of all deployments in the namespace as the target service.
- __true__: collect information about images for containers of deployments that __have versionExporter annotation__. 
Deployments that do not have the annotation will be ignored and not added to documentation.

