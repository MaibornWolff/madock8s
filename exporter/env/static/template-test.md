# ENV for service {{.Name}} at {{.Address}}
## This is a list of environment variables for deployment {{.DeploymentName}}<br/>
| Key | SourceType/Name.Key | Value |
|---|:-----------:|-----:|
{{range $key, $varArray := .ContainerVarsMap}}| __Container *`{{$key}}`*__|||<BR />
{{range $envVar := $varArray}}| {{$envVar.Key}} | {{if ne $envVar.ExtSourceType ""}}{{$envVar.ExtSourceType}}/{{$envVar.ExtSourceName}}.{{$envVar.ExtSourceKey}}{{end}} | {{$envVar.Value}} |<BR />
{{end}}{{end}}
