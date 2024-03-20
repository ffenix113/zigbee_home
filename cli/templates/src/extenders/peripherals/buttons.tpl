{{/* The templates are non-empty to force their usage. */}}
{{ define "top_level" }}
{{- range .Extender.Instances }}
static struct gpio_dt_spec {{.ID}} = GPIO_DT_SPEC_GET(DT_NODELABEL({{.ID}}), gpios);
{{- end }}
{{end}}

{{ define "loop"}} {{end}}

{{ define "main"}}
// Buttons will be configured by DK library.
{{end}}