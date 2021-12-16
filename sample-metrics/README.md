# Sample Metrics

Sample Metrics is a simple microservice, written in golang.
It exposes swagger-documented Rest API endpoints for demonstration.
The `/metrics` endpoint responsds with the list of available prometheus metrics and their descriptions.

## Exporters Configuration

The example configuration for available exporters is presented below. 
It is an excerpt from `yaml/deployment.yaml`.
Availability of `madock8s.exporter/<exporterName>` is sufficient for controller to notify the exporter with name <exporterName>.
For details on value meanings, please refer to corresponding exporter documentation.

```yaml
metadata:
  name: sample-metrics
  annotations:
    madock8s: SampleMetrics

    madock8s.exporter/prometheusExporter: true

    madock8s.exporter/gitlabExporter.baseurl: https://git.maibornwolff.de/api/v4/projects/3206/repository/
    madock8s.exporter/gitlabExporter.path: sample-metrics
    madock8s.exporter/gitlabExporter.recursive: true
    madock8s.exporter/gitlabExporter.ref: master
    madock8s.exporter/gitlabExporter.pattern: .\\.md

    madock8s.exporter/envExporter: sample-metrics

    madock8s.exporter/swaggerExporter.jsonurl: https://petstore.swagger.io/v2/swagger.json
    madock8s.exporter/swaggerExporter.baseurl: /docs
    madock8s.exporter/swaggerExporter.json: /swagger.json
    madock8s.exporter/swaggerExporter.port: 81
    madock8s.exporter/swaggerExporter: true

    madock8s.exporter/versionExporter: true
```


