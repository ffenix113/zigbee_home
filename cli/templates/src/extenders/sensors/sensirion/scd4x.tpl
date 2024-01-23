{{- define "top_level"}}
static const struct device *{{.Sensor.Label}}_{{.Endpoint}} = DEVICE_DT_GET(DT_NODELABEL({{.Sensor.Label}}));
{{- end}}

{{- define "main"}}
if (!{{.Sensor.Label}}_{{.Endpoint}}) {
    LOG_ERR("Failed to get {{.Sensor.Label}}");
    return ENODEV;
}
{{- end}}
