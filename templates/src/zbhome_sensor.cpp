#include <zephyr/logging/log.h>
#include <zephyr/drivers/sensor.h>
// TODO: maybe split this file to remove the header from here?
#include <zephyr/drivers/adc.h>

#include "zbhome_sensor.hpp"
#include "clusters.hpp"

LOG_MODULE_REGISTER(zbhome_sensor, LOG_LEVEL_DBG);

namespace zbhome
{
    namespace sensors
    {
        uint8_t setAttrValue(int endpoint, uint8_t *data_ptr, uint16_t clusterId, uint16_t valueId)
        {
            return (uint8_t)zb_zcl_set_attr_val(
                endpoint,
                clusterId,
                ZB_ZCL_CLUSTER_SERVER_ROLE,
                valueId,
                data_ptr,
                ZB_FALSE);
        }

        float convertSensorValue(struct sensor_value *value, float multiplier)
        {
            float measured_value = sensor_value_to_float(value);
            float calculated_value = (measured_value * multiplier);

            // LOG_DBG("Sensor converted value:\t%6f (raw: %d.%06d, mul: %f)", calculated_value, value->val1, value->val2, multiplier);

            return calculated_value;
        }

        int updateFetchedSamples(const struct device *sensor, int endpoint)
        {
            for (auto &config : sensorTypes)
            {
                struct sensor_value value;
                int err = sensor_channel_get(sensor, config.channel, &value);
                if (err)
                {
                    // If not supported - it is okay. We may try more channels that sensor defines.
                    if (err == -ENOTSUP)
                    {
                        continue;
                    }

                    LOG_ERR("Failed to get sensor %s channel %s: %d", sensor->name, config.channelName, err);
                    return err;
                }
                LOG_DBG("Sensor raw   %s/%s:\t%6d.%06d", sensor->name, config.channelName, value.val1, value.val2);

                float convertedValue = convertSensorValue(&value, config.multiplier);

                // CO2(and maybe some others) needs floating point precision, others can be slimmed to int16
                // TODO: review this code. Maybe it can be improved?
                union
                {
                    float f;
                    int16_t i;
                } val;

                if (config.channel == SENSOR_CHAN_CO2)
                {
                    val.f = convertedValue;
                }
                else
                {
                    val.i = (int16_t)convertedValue;
                }

                err = setAttrValue(endpoint, (zb_uint8_t *)&val, config.clusterId, config.attrValueId);
                if (err)
                {
                    LOG_ERR("Failed to set ZCL attribute for sensor %s, cluster %s: %d", sensor->name, config.channelName, err);
                    return err;
                }
            }

            return 0;
        }

        int read_adc_mv(const struct adc_dt_spec *spec, int32_t *valp)
        {
            int err;
            uint16_t buf;
            struct adc_sequence sequence = {
                .buffer = &buf,
                /* buffer size in bytes, not number of samples */
                .buffer_size = sizeof(buf),
            };

            (void)adc_sequence_init_dt(spec, &sequence);
            err = adc_read(spec->dev, &sequence);
            if (err < 0)
            {
                LOG_ERR("ADC %s@%d: Could not read (%d)\n", spec->dev->name, spec->channel_id, err);
                return err;
            }

            /*
             * If using differential mode, the 16 bit value
             * in the ADC sample buffer should be a signed 2's
             * complement value.
             */
            int32_t val_mv;
            if (spec->channel_cfg.differential)
            {
                val_mv = (int32_t)((int16_t)buf);
            }
            else
            {
                val_mv = (int32_t)buf;
            }

            LOG_DBG("ADC %s/%d raw value: %d", spec->dev->name, spec->channel_id, val_mv);
            err = adc_raw_to_millivolts_dt(spec, &val_mv);
            /* conversion to mV may not be supported, skip if not */
            if (err < 0)
            {
                LOG_DBG("  (value in mV not available)");
                return err;
            }

            LOG_DBG("ADC %s/%d mv value: %d", spec->dev->name, spec->channel_id, val_mv);

            *valp = val_mv;
            return 0;
        }
    }
}
