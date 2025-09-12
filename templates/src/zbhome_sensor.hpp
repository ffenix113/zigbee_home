#pragma once

#include <zephyr/drivers/sensor.h>

#include "types.hpp"
#include "clusters.hpp"

namespace zbhome
{
    namespace sensors
    {
        uint8_t setAttrValue(int endpoint, uint8_t *data_ptr, uint16_t clusterId, uint8_t valueId);
        int updateFetchedSamples(const struct device *sensor, int endpoint);
        int read_adc_mv(const struct adc_dt_spec *spec, int32_t *valp);

        struct sensorTypeConfig
        {
            sensor_channel channel;
            const char *channelName;
            uint16_t clusterId;
            uint8_t attrValueId;
            float multiplier;
        };

        // Maybe having full struct is not ideal here.
        // Reference or pointer may be better/cheaper?
        //
        // Names for channels are copied from Zephyr sensor shell code.
        // https://github.com/zephyrproject-rtos/zephyr/blob/72509c385e20410eebdc09d8ca60a206a063c064/drivers/sensor/sensor_shell.c#L48
        static const sensorTypeConfig sensorTypes[] = {
            {
                .channel = SENSOR_CHAN_AMBIENT_TEMP,
                .channelName = "ambient_temp",
                .clusterId = ZB_ZCL_CLUSTER_ID_TEMP_MEASUREMENT,
                .attrValueId = ZB_ZCL_ATTR_TEMP_MEASUREMENT_VALUE_ID,
                .multiplier = ZCL_TEMPERATURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER,
            },
            {
                .channel = SENSOR_CHAN_HUMIDITY,
                .channelName = "humidity",
                .clusterId = ZB_ZCL_CLUSTER_ID_REL_HUMIDITY_MEASUREMENT,
                .attrValueId = ZB_ZCL_ATTR_REL_HUMIDITY_MEASUREMENT_VALUE_ID,
                .multiplier = ZCL_HUMIDITY_MEASUREMENT_MEASURED_VALUE_MULTIPLIER,
            },
            {
                .channel = SENSOR_CHAN_PRESS,
                .channelName = "press",
                .clusterId = ZB_ZCL_CLUSTER_ID_PRESSURE_MEASUREMENT,
                .attrValueId = ZB_ZCL_ATTR_PRESSURE_MEASUREMENT_VALUE_ID,
                .multiplier = ZCL_PRESSURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER,
            },
// It is not always defined, so check before including.
#if ZB_ZCL_CLUSTER_ID_CARBON_DIOXIDE
            {
                .channel = SENSOR_CHAN_CO2,
                .channelName = "co2",
                .clusterId = ZB_ZCL_CLUSTER_ID_CARBON_DIOXIDE,
                .attrValueId = ZB_ZCL_ATTR_CARBON_DIOXIDE_VALUE_ID,
                .multiplier = ZCL_CARBON_DIOXIDE_MEASURED_VALUE_MULTIPLIER,
            },
#endif
        };
    }
}
