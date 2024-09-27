#pragma once

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

#define GENERATE_ATTR_SENSOR_VALUE_SETTER(cluster_name, measured_type, attribute_type, multiplier) \
    zb_zcl_status_t zbhome_set_attr_sensor_value_for_##cluster_name(int endpoint, struct sensor_value * value) { \
        measured_type measured_value = sensor_value_to_##measured_type(value); \
        attribute_type attr_value = (measured_value * multiplier); \
        return zbhome_set_attr_val_for_##cluster_name(endpoint, (zb_uint8_t*)&attr_value); \
    }

#define NAME_GENERATE_SENSOR_FULL_FOR_ATTR(cluster_name) int zbhome_sensor_fetch_and_update_##cluster_name(const struct device * sensor, int endpoint)
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

#ifdef ZB_ZCL_CLUSTER_ID_TEMP_MEASUREMENT
NAME_GENERATE_SENSOR_FULL_FOR_ATTR(temperature);
#endif

#ifdef ZB_ZCL_CLUSTER_ID_REL_HUMIDITY_MEASUREMENT
NAME_GENERATE_SENSOR_FULL_FOR_ATTR(humidity);
#endif

#ifdef ZB_ZCL_CLUSTER_ID_PRESSURE_MEASUREMENT
NAME_GENERATE_SENSOR_FULL_FOR_ATTR(pressure);
#endif

#ifdef ZB_ZCL_CLUSTER_ID_CARBON_DIOXIDE
NAME_GENERATE_SENSOR_FULL_FOR_ATTR(carbon_dioxide);
#endif
