{{- define "top_level"}}
static const struct device *{{.Sensor.Label}}_{{.Endpoint}} = DEVICE_DT_GET(DT_NODELABEL({{.Sensor.Label}}));
{{- end}}

{{- define "main"}}
if (!{{.Sensor.Label}}_{{.Endpoint}}) {
    LOG_ERR("Failed to get {{.Sensor.Label}}");
    return ENODEV;
}
{{- end}}

{{ /* 
    Yes, this loop is using bme680 logic,
    but it is not a problem, as they are basically identical,
    and I need to just make it work before refactoring to make it nice.
 */}}
{{- define "loop"}}
    int err = bosch_bme680_update_measurements({{.Sensor.Label}}_{{.Endpoint}});
	if (err) {
		LOG_ERR("Failed to upate measurements: %d", err);
	} else {
		err = bosch_bme680_update_temperature({{.Endpoint}}, {{.Sensor.Label}}_{{.Endpoint}});
		if (err) {
			LOG_ERR("Failed to update temperature: %d", err);
		}

		err = bosch_bme680_update_humidity({{.Endpoint}}, {{.Sensor.Label}}_{{.Endpoint}});
		if (err) {
			LOG_ERR("Failed to update humidity: %d", err);
		}
        {{ /* CO2 measurements are not supported by our Zigbee clusters yet. */ }}
	}

{{- end}}