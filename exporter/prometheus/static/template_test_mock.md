# Metrics from service {{.Name}} at {{.Address}}
{{range $b := .Blocks}}
**{{$b.Name}} as {{$b.Type}}**
**{{$b.Help}}**
{{end}}