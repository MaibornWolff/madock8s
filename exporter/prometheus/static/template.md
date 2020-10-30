# Metrics from service {{.Name}} at {{.Address}}
## This is a list of all prometheus-metrics for the given service<br/>
{{range $b := .Blocks}}[{{$b.Name}}]() / {{end}}
{{range $b := .Blocks}}
**{{$b.Type}}: {{$b.Name}}**<BR />
{{$b.Help}}
<BR />
{{end}}