# Swagger Exporter

Swagger Exporter processes swagger.json of a service to create a table of endpoints for the service. The table presents the following information: Endpoint, Method, Summary.

## Configuration

### Externally hosted Swagger configuration
Target deployment has to provide the following annotations in metadata:

```yaml
metadata:
  annotations:
    madock8s.exporter/swaggerExporter.jsonurl: https://petstore.swagger.io/v2/swagger.json
```

Breakdown:
- jsonurl: The external URL of swagger.json. If jsonurl is set, other parameters are ignored.

### Internally hosted Swagger configuration
Target deployment has to provide the following annotations in metadata:

```yaml
metadata:
  annotations:
    madock8s.exporter/swaggerExporter.port: 81
    madock8s.exporter/swaggerExporter.baseurl: /docs
    madock8s.exporter/swaggerExporter.json: /swagger.json
```

Breakdown:
- port: Host port of Swagger endpoint.
- baseurl: The path component added to service's address when accessing swagger.json
- json: Filename of Swagger OpenAPI2.0 configuration. Defaults to `swagger.json`

The exporter combines cluster-ip (e.g, 10.20.170.89), port, baseurl, and json params to fetch the file. 
For example, the final URL to access internally hosted swagger config is: 
```
  10.20.170.89:81/docs/swagger.json
```

### Almost 0-Config

It is possible to pre-configure Swagger Exporter using envirionment variables. 
Prepequisite: most of the services have same Swagger configuration (baseurl, json and port).
If Swagger configuration is different for a deployment that needs to be exported, add required annotations as described in steps above.
Configuration from annotations has higher priority.

1. To enable Swagger Exporter for deployment, add the following annotation to its yaml-definition.
```yaml
metadata:
  annotations:
    madock8s.exporter/swaggerExporter: true
```

2. Env Variables for Swagger Exporter can be set in helm /swagger-exporter/values.yaml in fields of `env` group.
Note: do NOT change field names.

Exempt from values.yaml:
```yaml
env:
  # madock8s.exporter/swaggerExporter.jsonurl
  swaggerJsonURL: https://petstore.swagger.io/v2/swagger.json
  # madock8s.exporter/swaggerExporter.port
  swaggerPort: 81
  # madock8s.exporter/swaggerExporter.baseurl
  baseURL: /docs
  # madock8s.exporter/swaggerExporter.json
  swaggerJSON: /swagger.json
```

### Delete Mode

Configure behaviour for md-files when target deployment is deleted.
Simply update helm/values.yaml/env.deletionStrategy to one of the following values:
- IGNORE - keep existing md-files without modification;
- UPDATE - add "Deployments of the service were deleted on timestamp" to existing md-file (default);
- DELETE - remove the respective md-file from storage.
