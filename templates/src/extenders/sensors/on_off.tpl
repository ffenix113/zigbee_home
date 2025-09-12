{{ define "top_level" }} {{ end }}
{{ define "loop" }} {{end}}
{{ define "main"}} {{ end}}

{{ define "component_constructor_arguments"}}
const struct gpio_dt_spec {{.Sensor.Pin.Label}} = GPIO_DT_SPEC_GET(DT_NODELABEL({{.Sensor.Pin.Label}}), gpios);
if (!gpio_is_ready_dt(&{{.Sensor.Pin.Label}})) {
    LOG_ERR("Pin {{.Sensor.Pin.Label}} is not ready");
    return -1;
}

int err = gpio_pin_configure_dt(&{{.Sensor.Pin.Label}}, GPIO_OUTPUT);
if (err != 0) {
    LOG_ERR("Cannot configure pin {{.Sensor.Pin.Label}}");
    return err;
}
{{end}}

{{ define "component_constructor_argument_names"}} {{.Sensor.Pin.Label}} {{end}}