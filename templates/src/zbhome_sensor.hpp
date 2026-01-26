#pragma once

#include <zephyr/drivers/sensor.h>

#ifdef __cplusplus
extern "C"
{
#endif
// Not ideal, maybe there is a better way.
#include <zboss_api.h>
#include <zcl/zb_zcl_el_measurement.h>

#ifdef __cplusplus
}
#endif

#include "types.hpp"
#include "clusters.hpp"

namespace zbhome
{
    namespace sensors
    {
        uint8_t setAttrValue(int endpoint, uint8_t *data_ptr, uint16_t clusterId, uint16_t valueId);
        int updateFetchedSamples(const struct device *sensor, int endpoint);
        int read_adc_mv(const struct adc_dt_spec *spec, int32_t *valp);

        struct sensorTypeConfig
        {
            sensor_channel channel;
            const char *channelName;
            uint16_t clusterId;
            uint16_t attrValueId;
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
#ifdef ZB_ZCL_CLUSTER_ID_CARBON_DIOXIDE
            {
                .channel = SENSOR_CHAN_CO2,
                .channelName = "co2",
                .clusterId = ZB_ZCL_CLUSTER_ID_CARBON_DIOXIDE,
                .attrValueId = ZB_ZCL_ATTR_CARBON_DIOXIDE_VALUE_ID,
                .multiplier = ZCL_CARBON_DIOXIDE_MEASURED_VALUE_MULTIPLIER,
            },
#endif
#ifdef ZCL_DC_VOLTAGE_VALUE_DIVISOR
            {
                .channel = SENSOR_CHAN_VOLTAGE,
                .channelName = "dc_voltage",
                .clusterId = ZB_ZCL_CLUSTER_ID_ELECTRICAL_MEASUREMENT,
                .attrValueId = ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DC_VOLTAGE_ID,
                .multiplier = ZCL_DC_VOLTAGE_VALUE_DIVISOR,
            },
            {
                .channel = SENSOR_CHAN_CURRENT,
                .channelName = "dc_current",
                .clusterId = ZB_ZCL_CLUSTER_ID_ELECTRICAL_MEASUREMENT,
                .attrValueId = ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DC_CURRENT_ID,
                .multiplier = ZCL_DC_CURRENT_VALUE_DIVISOR,
            },
            {
                .channel = SENSOR_CHAN_POWER,
                .channelName = "dc_power",
                .clusterId = ZB_ZCL_CLUSTER_ID_ELECTRICAL_MEASUREMENT,
                .attrValueId = ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DCPOWER_ID,
                .multiplier = ZCL_DC_POWER_VALUE_DIVISOR,
            },
#endif
        };
    }
}
