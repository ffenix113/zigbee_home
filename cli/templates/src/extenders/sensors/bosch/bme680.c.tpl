/*
 * Copyright (c) 2022 Nordic Semiconductor ASA
 *
 * SPDX-License-Identifier: LicenseRef-Nordic-5-Clause
 */

#include <zephyr/drivers/sensor.h>
#include <zephyr/logging/log.h>

#include "bosch_bme680.h"

LOG_MODULE_DECLARE(app, LOG_LEVEL_INF);

/*
 * Sensor value is represented as having an integer and a fractional part,
 * and can be obtained using the formula val1 + val2 * 10^(-6).
 */
// #define BOSCH_BME680_DIVISOR 1000000

// /*
//  * Sensor value is represented as having an integer and a fractional part,
//  * and can be obtained using the formula val1 + val2 * 10^(-6). Negative
//  * values also adhere to the above formula, but may need special attention.
//  * Here are some examples of the value representation:
//  *
//  *      0.5: val1 =  0, val2 =  500000
//  *     -0.5: val1 =  0, val2 = -500000
//  *     -1.0: val1 = -1, val2 =  0
//  *     -1.5: val1 = -1, val2 = -500000
//  */
// static float convert_bosch_bme680_value(struct sensor_value value)
// {
// 	float result = 0.0f;

// 	/* Determine sign */
// 	result = (value.val1 < 0 || value.val2 < 0) ? -1.0f : 1.0f;

// 	/* Use absolute values */
// 	value.val1 = value.val1 < 0 ? -value.val1 : value.val1;
// 	value.val2 = value.val2 < 0 ? -value.val2 : value.val2;

// 	/* Calculate value */
// 	result *= (value.val1 + value.val2 / (float)BOSCH_BME680_DIVISOR);

// 	return result;
// }

int bosch_bme680_update_measurements(const struct device *sensor)
{
	int err = 0;
	err = sensor_sample_fetch(sensor);

	return err;
}

int bosch_bme680_get_temperature(const struct device *sensor, float *temperature)
{
	int err = 0;

	struct sensor_value sensor_temperature;

	err = sensor_channel_get(sensor,
					SENSOR_CHAN_AMBIENT_TEMP,
					&sensor_temperature);
	if (err) {
		LOG_ERR("Failed to get sensor channel: %d", err);
	} else {
		LOG_INF("Sensor    T:%3d.%06d [*C]",
			sensor_temperature.val1, sensor_temperature.val2);
		*temperature = sensor_value_to_float(&sensor_temperature);
	}

	return err;
}

int bosch_bme680_get_pressure(const struct device *sensor, float *pressure)
{
	int err = 0;

	struct sensor_value sensor_pressure;

	err = sensor_channel_get(sensor,
					SENSOR_CHAN_PRESS,
					&sensor_pressure);
	if (err) {
		LOG_ERR("Failed to get sensor channel: %d", err);
	} else {
		LOG_INF("Sensor    P:%3d.%06d [kPa]",
			sensor_pressure.val1, sensor_pressure.val2);
		*pressure = sensor_value_to_float(&sensor_pressure);
	}

	return err;
}

int bosch_bme680_get_humidity(const struct device *sensor, float *humidity)
{
	int err = 0;

	struct sensor_value sensor_humidity;

	err = sensor_channel_get(sensor,
					SENSOR_CHAN_HUMIDITY,
					&sensor_humidity);
	if (err) {
		LOG_ERR("Failed to get sensor channel: %d", err);
	} else {
		LOG_INF("Sensor    H:%3d.%06d [%%]",
			sensor_humidity.val1, sensor_humidity.val2);
		*humidity = sensor_value_to_float(&sensor_humidity);
	}

	return err;
}

int bosch_bme680_update_temperature(int endpoint, const struct device *sensor)
{
	int err = 0;

	float measured_temperature = 0.0f;
	int16_t temperature_attribute = 0;

	err = bosch_bme680_get_temperature(sensor, &measured_temperature);
	if (err) {
		LOG_ERR("Failed to get sensor temperature: %d", err);
	} else {
		/* Convert measured value to attribute value, as specified in ZCL */
		temperature_attribute = (int16_t)
					(measured_temperature *
					 ZCL_TEMPERATURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER);
		LOG_INF("Attribute T:%10d", temperature_attribute);

		/* Set ZCL attribute */
		zb_zcl_status_t status = zb_zcl_set_attr_val(endpoint,
							     ZB_ZCL_CLUSTER_ID_TEMP_MEASUREMENT,
							     ZB_ZCL_CLUSTER_SERVER_ROLE,
							     ZB_ZCL_ATTR_TEMP_MEASUREMENT_VALUE_ID,
							     (zb_uint8_t *)&temperature_attribute,
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

int bosch_bme680_update_humidity(int endpoint, const struct device *sensor)
{
	int err = 0;

	float measured_humidity = 0.0f;
	int16_t humidity_attribute = 0;

	err = bosch_bme680_get_humidity(sensor, &measured_humidity);
	if (err) {
		LOG_ERR("Failed to get sensor humidity: %d", err);
	} else {
		/* Convert measured value to attribute value, as specified in ZCL */
		humidity_attribute = (int16_t)
				     (measured_humidity *
				      ZCL_HUMIDITY_MEASUREMENT_MEASURED_VALUE_MULTIPLIER);
		LOG_INF("Attribute H:%10d", humidity_attribute);

		zb_zcl_status_t status = zb_zcl_set_attr_val(
			endpoint,
			ZB_ZCL_CLUSTER_ID_REL_HUMIDITY_MEASUREMENT,
			ZB_ZCL_CLUSTER_SERVER_ROLE,
			ZB_ZCL_ATTR_REL_HUMIDITY_MEASUREMENT_VALUE_ID,
			(zb_uint8_t *)&humidity_attribute,
			ZB_FALSE);
		if (status) {
			LOG_ERR("Failed to set ZCL attribute: %d", status);
			err = status;
		}
	}

	return err;
}
