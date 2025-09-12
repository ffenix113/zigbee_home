{{/* The templates are non-empty to force their usage. */}}
{{ define "top_level" }} {{end}}
{{ define "button_changed"}} {{end}}
{{ define "loop"}} {{end}}
{{ define "main"}} {{end}}

{{ define "component_constructor_arguments"}}
const struct adc_dt_spec adc_channel_{{.Sensor.ADCPin.Name}} = ADC_DT_SPEC_GET_BY_IDX(DT_PATH(zephyr_user), 0);
if (int err = adc_channel_setup_dt(&adc_channel_{{.Sensor.ADCPin.Name}}); err < 0) {
    LOG_ERR("Could not setup adc channel '{{.Sensor.ADCPin.Name}}': %d\n", err);
    return false;
}
{{end}}
{{ define "component_constructor_argument_names"}} adc_channel_{{.Sensor.ADCPin.Name}}, {{.Sensor.MaxMoistureMv}}, {{.Sensor.MinMoistureMv}} {{end}}