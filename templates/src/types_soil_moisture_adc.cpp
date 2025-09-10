#include <stdint.h>

#include <zephyr/drivers/adc.h>
#include <zephyr/logging/log.h>

#include "zbhome_sensor.hpp"
#include "types_soil_moisture_adc.hpp"

LOG_MODULE_REGISTER(soil_moisture_adc, LOG_LEVEL_INF);

namespace zbhome
{
    namespace components
    {
        void SoilMoistureADC::onLoop()
        {
            uint16_t mv_val;
            int err = zbhome::sensors::read_adc_mv(&m_adc_spec, &mv_val);
            if (err < 0)
            {
                LOG_ERR("failed to fetch adc mv value: %d", err);
                return;
            }

            setAttrValue(getEndpoint(), (uint8_t *)(&mv_val), (0x0408), ZB_ZCL_ATTR_REL_HUMIDITY_MEASUREMENT_VALUE_ID);
        };
    }
}