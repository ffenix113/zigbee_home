#include <stdint.h>

#include <zephyr/drivers/adc.h>
#include <zephyr/logging/log.h>

#include "zbhome_sensor.hpp"
#include "types_power_config.hpp"

LOG_MODULE_REGISTER(power_config, LOG_LEVEL_DBG);

namespace zbhome
{
    namespace components
    {
        void PowerConfig::onLoop() {
            uint16_t batt_mv;
            int err = zbhome::sensors::read_adc_mv(&m_adc_spec, &batt_mv);
            if (err < 0) {
                LOG_ERR("failed to fetch adc mv value: %d", err);
                return;
            }

            // Set the voltage
            zbhome::sensors::setAttrValue(getEndpoint(), (uint8_t *)(&batt_mv), (0x0001), ZB_ZCL_ATTR_POWER_CONFIG_BATTERY_VOLTAGE_ID);

            // From ZCL specification, 3.3.2.2.3.2
            // "a range between zero and 100%, with 0x00 = 0%, 0x64 = 50%, and 0xC8 = 100%".
            // 0xC8 == 200.
            uint8_t percentage; 
            if (batt_mv >= m_rated_mv) {
                percentage = 200;
            } else if (batt_mv <= m_min_mv_threshold) {
                percentage = 0;
            } else {
                percentage = (((batt_mv-m_min_mv_threshold)*200) / (m_rated_mv - m_min_mv_threshold));
            }

            // Have only 10% increments
            percentage = ((percentage + 10) / 20) * 20;

            LOG_DBG("ADC power config {{.Sensor.ADCPin.Name}}: rated: %d, threshold: %d, current mv: %d, percent: %d", m_rated_mv, m_min_mv_threshold, batt_mv, percentage/2);

            // Set percentage
            zbhome::sensors::setAttrValue(getEndpoint(), (uint8_t*)(&percentage), (0x0001), ZB_ZCL_ATTR_POWER_CONFIG_BATTERY_PERCENTAGE_REMAINING_ID);
        };
    }
}