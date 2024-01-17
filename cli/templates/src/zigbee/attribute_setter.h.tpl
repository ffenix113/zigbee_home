#define ATTR_SET_VAL_FOR_TYPE(attr_type) \
    int attr_set_val_##attr_type##(int endpoint, attr_type value) { \
    
    }    

int attr_set_val_int(int endpoint, const struct device *sensor)
{
	int err = 0;

	float measured_pressure = 0.0f;
	int16_t pressure_attribute = 0;

	err = bosch_bme680_get_pressure(sensor, &measured_pressure);
	if (err) {
		LOG_ERR("Failed to get sensor pressure: %d", err);
	} else {
		/* Convert measured value to attribute value, as specified in ZCL */
		pressure_attribute = (int16_t)
				     (measured_pressure *
				      ZCL_PRESSURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER);
		LOG_INF("Attribute P:%10d", pressure_attribute);

		/* Set ZCL attribute */
		zb_zcl_status_t status = zb_zcl_set_attr_val(
			endpoint,
			ZB_ZCL_CLUSTER_ID_PRESSURE_MEASUREMENT,
			ZB_ZCL_CLUSTER_SERVER_ROLE,
			ZB_ZCL_ATTR_PRESSURE_MEASUREMENT_VALUE_ID,
			(zb_uint8_t *)&pressure_attribute,
			ZB_FALSE);
		if (status) {
			LOG_ERR("Failed to set ZCL attribute: %d", status);
			err = status;
		}
	}

	return err;
}

int bosch_bme680_update_pressure(int endpoint, const struct device *sensor)
{
	int err = 0;

	float measured_pressure = 0.0f;
	int16_t pressure_attribute = 0;

	err = bosch_bme680_get_pressure(sensor, &measured_pressure);
	if (err) {
		LOG_ERR("Failed to get sensor pressure: %d", err);
	} else {
		/* Convert measured value to attribute value, as specified in ZCL */
		pressure_attribute = (int16_t)
				     (measured_pressure *
				      ZCL_PRESSURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER);
		LOG_INF("Attribute P:%10d", pressure_attribute);

		/* Set ZCL attribute */
		zb_zcl_status_t status = zb_zcl_set_attr_val(
			endpoint,
			ZB_ZCL_CLUSTER_ID_PRESSURE_MEASUREMENT,
			ZB_ZCL_CLUSTER_SERVER_ROLE,
			ZB_ZCL_ATTR_PRESSURE_MEASUREMENT_VALUE_ID,
			(zb_uint8_t *)&pressure_attribute,
			ZB_FALSE);
		if (status) {
			LOG_ERR("Failed to set ZCL attribute: %d", status);
			err = status;
		}
	}

	return err;
}