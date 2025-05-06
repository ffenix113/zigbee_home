#include <zephyr/logging/log.h>
#include <zephyr/drivers/sensor.h>

#include "zbhome_sensor.hpp"
#include "clusters.hpp"

LOG_MODULE_DECLARE(app, LOG_LEVEL_INF);

namespace zbhome {
    namespace sensors {
        zb_zcl_status_t setAttrValue(int endpoint, zb_uint8_t * data_ptr, uint16_t clusterId, uint8_t valueId) {
            return zb_zcl_set_attr_val(
                endpoint,
                clusterId,
                ZB_ZCL_CLUSTER_SERVER_ROLE,
                valueId,
                data_ptr,
                ZB_FALSE);
        }
        
        uint8_t convertAttrValue(struct sensor_value * value, uint8_t multiplier) {
            float measured_value = sensor_value_to_float(value);
            return (uint8_t)(measured_value * multiplier);
        }
        
        int updateFetchedSamples(const struct device * sensor, int endpoint) {
            for (auto& config : sensorTypes) {
                struct sensor_value value;
                int err = sensor_channel_get(sensor, config.channel, &value);
                if (err) {
                    // If not supported - it is okay. We may try more channels that sensor defines.
                    if (err == -ENOTSUP) {
                        continue;
                    }

                    LOG_ERR("Failed to get sensor %s channel %s: %d", sensor->name, config.channelName, err);
                    return err;
                }
                LOG_DBG("Sensor raw   %s/%s:\t%6d.%06d", sensor->name, config.channelName, value.val1, value.val2);

                auto convertedValue = convertAttrValue(&value, config.multiplier);

                err = setAttrValue(endpoint, (zb_uint8_t*)&convertedValue, config.clusterId, config.attrValueId);
                if (err) {
                    LOG_ERR("Failed to set ZCL attribute for sensor %s, cluster %s: %d", sensor->name, config.channelName, err);
                    return err;
                }
            }

            return 0;
        }
    }
}
