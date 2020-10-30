# Image versions for deployments from namespace {{ .Namespace }}

| Deployment | Container | Image | Version |
| :--------- | :-----------: | :-------: | -----------: |
{{ range $deployment := .Deployments }}{{ range $c := $deployment.Containers }}| {{ $deployment.Name }} | {{ $c.Name }} | {{ $c.Image }}  | {{ $c.Version }} |
{{ end }}{{ end }}
