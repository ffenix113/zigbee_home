#ifndef BOSCH_BME680_H
#define BOSCH_BME680_H

#include <zboss_api.h>
#include <zboss_api_addons.h>
#include <zcl/zb_zcl_temp_measurement_addons.h>

/* Zigbee Cluster Library 4.4.2.2.1.1: MeasuredValue = 100x temperature in degrees Celsius */
#define ZCL_TEMPERATURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER 100
/* Zigbee Cluster Library 4.5.2.2.1.1: MeasuredValue = 10x pressure in kPa */
#define ZCL_PRESSURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER 10
/* Zigbee Cluster Library 4.7.2.1.1: MeasuredValue = 100x water content in % */
#define ZCL_HUMIDITY_MEASUREMENT_MEASURED_VALUE_MULTIPLIER 100

/**
 * @brief Updates internal measurements performed by sensor.
 *
 * @note It has to be called each time a fresh measurements are required.
 *	 It does not change any ZCL attributes.
 *
 * @return 0 if success, error code if failure.
 */
int bosch_bme680_update_measurements(const struct device *sensor);

/**
 * @brief Updates ZCL temperature attribute using value obtained during last weather check.
 *
 * @return 0 if success, error code if failure.
 */
int bosch_bme680_update_temperature(int endpoint, const struct device *sensor);

/**
 * @brief Updates ZCL pressure attribute using value obtained during last weather check.
 *
 * @return 0 if success, error code if failure.
 */
int bosch_bme680_update_pressure(int endpoint, const struct device *sensor);

/**
 * @brief Updates ZCL relative humidity attribute using value obtained during last weather check.
 *
 * @return 0 if success, error code if failure.
 */
int bosch_bme680_update_humidity(int endpoint, const struct device *sensor);

#endif