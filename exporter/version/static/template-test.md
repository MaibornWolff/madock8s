# Image versions for deployments from namespace {{ .Namespace }}

| Deployment | ContainerName | ImageName | ImageVersion |
| :--------- | :-----------: | :-------: | -----------: |
{{ range $deployment := .Deployments }}{{ range $c := $deployment.Containers }}| {{ $deployment.Name }} | {{ $c.Name }} | {{ $c.Image }} | {{ $c.Version }} |
{{ end }}{{ end }}
