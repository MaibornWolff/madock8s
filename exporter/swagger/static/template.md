# Swagger-documented endpoints for service {{ .Name }} at {{ .Address }}

- Swagger.json location: {{ .JSONURL }}
- Swagger home page (internal access only): {{ .BaseURL }} 

## List of available endpoints:

| Path | Method | Summary |
| :---: | :----: | ------: |
{{ range $path, $methods := .Endpoints }}{{ range $methods }} | `{{ $path }}` | {{ .Method }} | {{ .Op.Summary }} |
{{ end }}{{ end }}
