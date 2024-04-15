{{define "top_level"}}
/* Data of ADC io-channels specified in devicetree. */
{{ range $i, $instance := .Extender.Instances }}
static const struct adc_dt_spec adc_channel_{{$instance.Name}} = ADC_DT_SPEC_GET_BY_IDX(DT_PATH(zephyr_user), {{$i}});
{{end}}
{{end}}

{{ define "loop"}} {{end}}


{{ define "main"}}
int err;
{{ range .Extender.Instances }}
err = adc_channel_setup_dt(&adc_channel_{{.Name}});
if (err < 0) {
    LOG_ERR("Could not setup channel '{{.Name}}' (%d)\n", err);
    return 0;
}
{{end}}
{{end}}