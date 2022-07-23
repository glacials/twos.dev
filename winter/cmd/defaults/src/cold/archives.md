---
type: page
---

# Archives

{{ if archives }}
{{$prev := (index archives 0).Year}}
{{range archives}}

{{if lt .Year (sub $prev 1) }}

.

.

.

{{end}}

## {{.Year}}

{{range .Documents}}

- {{with .Category}}{{.}}:{{end}} [{{.Title}}]({{.Shortname}}.html)

{{end}}
{{$prev = .Year}}
{{end}}
{{end}}
