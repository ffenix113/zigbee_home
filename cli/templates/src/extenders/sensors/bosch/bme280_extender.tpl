{{- define "top_level"}}
static const struct device *{{.Sensor.Label}}_{{.Endpoint}};
{{- end}}

{{- define "main"}}
{{.Sensor.Label}}_{{.Endpoint}} = DEVICE_DT_GET(DT_NODELABEL({{.Sensor.Label}}));
if (!{{.Sensor.Label}}_{{.Endpoint}}) {
    LOG_ERR("Failed to get {{.Sensor.Label}}");
    return ENODEV;
}
{{- end}}

{{- define "loop"}}
    int err = bosch_bme680_update_measurements({{.Sensor.Label}}_{{.Endpoint}});
	if (err) {
		LOG_ERR("Failed to upate measurements: %d", err);
	} else {
		err = bosch_bme680_update_temperature({{.Endpoint}}, {{.Sensor.Label}}_{{.Endpoint}});
		if (err) {
			LOG_ERR("Failed to update temperature: %d", err);
		}

		err = bosch_bme680_update_pressure({{.Endpoint}}, {{.Sensor.Label}}_{{.Endpoint}});
		if (err) {
			LOG_ERR("Failed to update pressure: %d", err);
		}

		err = bosch_bme680_update_humidity({{.Endpoint}}, {{.Sensor.Label}}_{{.Endpoint}});
		if (err) {
			LOG_ERR("Failed to update humidity: %d", err);
		}
	}

{{- end}}