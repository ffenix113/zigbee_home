#define GENERATE_ATTR_VAL_SETTER(cluster_name, cluster) \
    zb_zcl_status_t zbhome_set_attr_val_for_##cluster_name(int endpoint, zb_uint8_t * data_ptr) { \
        return zb_zcl_set_attr_val( \
            endpoint, \
            ZB_ZCL_CLUSTER_ID_##cluster, \
            ZB_ZCL_CLUSTER_SERVER_ROLE, \
            ZB_ZCL_ATTR_##cluster##_VALUE_ID, \
            data_ptr, \
            ZB_FALSE); \
    }


GENERATE_ATTR_VAL_SETTER(temperature, TEMP_MEASUREMENT)
GENERATE_ATTR_VAL_SETTER(humidity, REL_HUMIDITY_MEASUREMENT)
GENERATE_ATTR_VAL_SETTER(pressure, PRESSURE_MEASUREMENT)
GENERATE_ATTR_VAL_SETTER(carbon_dioxide, CARBON_DIOXIDE)

#define GENERATE_ATTR_SENSOR_VALUE_SETTER(cluster_name, measured_type, attribute_type, multiplier) \
    zb_zcl_status_t zbhome_set_attr_sensor_value_for_##cluster_name(int endpoint, struct sensor_value * value) { \
        measured_type measured_value = sensor_value_to_##measured_type(value); \
        attribute_type attr_value = (measured_value * multiplier); \
        return zbhome_set_attr_val_for_##cluster_name(endpoint, (zb_uint8_t*)&attr_value); \
    }



GENERATE_ATTR_SENSOR_VALUE_SETTER(temperature, float, int16_t, ZCL_TEMPERATURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER)
GENERATE_ATTR_SENSOR_VALUE_SETTER(humidity, float, int16_t, ZCL_HUMIDITY_MEASUREMENT_MEASURED_VALUE_MULTIPLIER)
GENERATE_ATTR_SENSOR_VALUE_SETTER(pressure, float, int16_t, ZCL_PRESSURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER)
GENERATE_ATTR_SENSOR_VALUE_SETTER(carbon_dioxide, float, float, ZCL_CARBON_DIOXIDE_MEASURED_VALUE_MULTIPLIER)


#define GENERATE_SENSOR_FULL_FOR_ATTR(cluster_name, sensor_channel) \
    int zbhome_sensor_fetch_and_update_##cluster_name(const struct device * sensor, int endpoint) { \
		struct sensor_value value; \
		int err = sensor_channel_get(sensor, sensor_channel, &value); \
        if (err) { \
            LOG_ERR("Failed to get sensor %s channel %s: %d", sensor->name, #cluster_name, err); \
            return err; \
        } \
        LOG_INF("Sensor raw   %s/%s: %6d.%06d", sensor->name, #cluster_name, value.val1, value.val2); \
        err = (int)zbhome_set_attr_sensor_value_for_##cluster_name(endpoint, &value); \
        if (err) { \
            LOG_ERR("Failed to set ZCL attribute for sensor %s, cluster %s: %d", sensor->name, #cluster_name, err); \
        } \
        return err; \
    }

GENERATE_SENSOR_FULL_FOR_ATTR(temperature, SENSOR_CHAN_AMBIENT_TEMP)
GENERATE_SENSOR_FULL_FOR_ATTR(humidity, SENSOR_CHAN_HUMIDITY)
GENERATE_SENSOR_FULL_FOR_ATTR(pressure, SENSOR_CHAN_PRESS)
GENERATE_SENSOR_FULL_FOR_ATTR(carbon_dioxide, SENSOR_CHAN_CO2)