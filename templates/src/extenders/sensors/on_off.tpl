{{ define "top_level" }}
static const struct gpio_dt_spec {{.Sensor.Pin.Label}} = GPIO_DT_SPEC_GET(DT_NODELABEL({{.Sensor.Pin.Label}}), gpios);
{{ end }}

{{ define "loop" }} {{end}}

{{ define "main"}}
	if (!gpio_is_ready_dt(&{{.Sensor.Pin.Label}})) {
        LOG_ERR("Pin {{.Sensor.Pin.Label}} is not ready");
		return -1;
	}

    int err = gpio_pin_configure_dt(&{{.Sensor.Pin.Label}}, GPIO_INPUT);
    if (err != 0) {
        LOG_ERR("Cannot configure pin {{.Sensor.Pin.Label}}");
        return err;
    }
{{ end}}