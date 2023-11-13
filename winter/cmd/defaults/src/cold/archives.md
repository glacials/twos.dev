---
type: page
---

# Archives

{{ if yearly posts }}
{{$prev := (index (yearly posts) 0).Year}}
{{range (yearly posts)}}

{{if lt .Year (sub $prev 1) }}

.

.

.

{{end}}

## {{.Year}}

{{range .Documents}}

- {{with .Category}}{{.}}:{{end}} [{{.Title}}]({{.WebPath}})

{{end}}
{{$prev = .Year}}
{{end}}
{{end}}
