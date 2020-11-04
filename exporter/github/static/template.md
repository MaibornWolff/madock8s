# Markdown files for service {{.Name}} at {{.Address}}
## These are all markdown files found in repository of the service


{{range $idx, $file := .MdFiles}}

`{{$file.Path}}`
---
{{$file.Content}}

---


{{end}}
