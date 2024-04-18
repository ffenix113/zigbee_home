{{/* The templates are non-empty to force their usage. */}}
{{ define "top_level" }}
{{- range .Extender.Instances }}
static struct gpio_dt_spec {{.ID}} = GPIO_DT_SPEC_GET(DT_NODELABEL({{.ID}}), gpios);
{{- end }}
{{end}}

{{ define "loop"}}

{{end}}

{{ define "main"}}
int ret;
{{- range .Extender.Instances}}
ret = gpio_pin_configure_dt(&{{.ID}}, GPIO_OUTPUT);
if (ret != 0) {
    LOG_ERR("Error %d: failed to configure LED device %s pin %d\n",
            ret, {{.ID}}.port->name, {{.ID}}.pin);
    return ret;
}
{{- end}}
{{end}}