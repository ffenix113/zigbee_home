{{- define "top_level"}}
static const struct device *{{.Sensor.Label}}_{{.Endpoint}} = DEVICE_DT_GET(DT_NODELABEL({{.Sensor.Label}}));
{{- end}}

{{- define "main"}}
if (!{{.Sensor.Label}}_{{.Endpoint}}) {
    LOG_ERR("Failed to get {{.Sensor.Label}}");
    return ENODEV;
}
{{- end}}

{{/* 
    Yes, this loop is using bme680 logic,
    but it is not a problem, as they are basically identical,
    and I need to just make it work before refactoring to make it nice.
 */}}
{{- define "loop"}}
    int err = sensor_sample_fetch({{.Sensor.Label}}_{{.Endpoint}});
	if (err) {
		LOG_ERR("Failed to upate {{.Sensor.Label}} measurements: %d", err);
	} else {
		struct sensor_value temp, hum, co2;
		float temp_v, hum_v, co2_v;

		sensor_channel_get({{.Sensor.Label}}_{{.Endpoint}}, SENSOR_CHAN_AMBIENT_TEMP, &temp);
		sensor_channel_get({{.Sensor.Label}}_{{.Endpoint}}, SENSOR_CHAN_HUMIDITY, &hum);
		sensor_channel_get({{.Sensor.Label}}_{{.Endpoint}}, SENSOR_CHAN_CO2, &co2);

		temp_v = sensor_value_to_float(&temp);
		hum_v = sensor_value_to_float(&hum);
		co2_v = sensor_value_to_float(&co2);

		LOG_INF("Sensor    SCD4X/T:%3d.%06d [*C]", temp.val1, temp.val2);
		int16_t temperature_attribute = (int16_t)
					(temp_v *
					 ZCL_TEMPERATURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER);
		LOG_INF("Attribute SCD4X/T:%10d", temperature_attribute);

		/* Set ZCL attribute */
		zb_zcl_status_t status = zb_zcl_set_attr_val({{.Endpoint}},
							     ZB_ZCL_CLUSTER_ID_TEMP_MEASUREMENT,
							     ZB_ZCL_CLUSTER_SERVER_ROLE,
							     ZB_ZCL_ATTR_TEMP_MEASUREMENT_VALUE_ID,
							     (zb_uint8_t *)&temperature_attribute,
							     ZB_FALSE);
		if (status) {
			LOG_ERR("Failed to set ZCL attribute: %d", status);
			err = status;
		}


		LOG_INF("Sensor    SCD4X/H:%3d.%06d [%%]", hum.val1, hum.val2);
		int16_t humidity_attribute = (int16_t)
					(hum_v *
					 ZCL_HUMIDITY_MEASUREMENT_MEASURED_VALUE_MULTIPLIER);
		LOG_INF("Attribute SCD4X/H:%10d", humidity_attribute);

		/* Set ZCL attribute */
		status = zb_zcl_set_attr_val({{.Endpoint}},
							     ZB_ZCL_CLUSTER_ID_REL_HUMIDITY_MEASUREMENT,
							     ZB_ZCL_CLUSTER_SERVER_ROLE,
							     ZB_ZCL_ATTR_REL_HUMIDITY_MEASUREMENT_VALUE_ID,
							     (zb_uint8_t *)&humidity_attribute,
							     ZB_FALSE);
		if (status) {
			LOG_ERR("Failed to set ZCL attribute: %d", status);
			err = status;
		}

		LOG_INF("Sensor    SCD4X/CO2:%4d [ppm]", co2.val1);

		/* Convert measured value to attribute value, as specified in ZCL */
		float co2_attribute = co2_v * ZCL_CARBON_DIOXIDE_MEASURED_VALUE_MULTIPLIER;
		LOG_INF("Attribute SCD4X/CO2:%f", co2_attribute);

		/* Set ZCL attribute */
		status = zb_zcl_set_attr_val(
			{{.Endpoint}},
			ZB_ZCL_CLUSTER_ID_CARBON_DIOXIDE,
			ZB_ZCL_CLUSTER_SERVER_ROLE,
			ZB_ZCL_ATTR_CARBON_DIOXIDE_VALUE_ID,
			(zb_uint8_t *)&co2_attribute,
			ZB_FALSE);
		if (status) {
			LOG_ERR("Failed to set ZCL attribute: %d", status);
			err = status;
		}
	}

{{- end}}