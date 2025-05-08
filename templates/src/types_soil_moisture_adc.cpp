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
        void SoilMoistureADC::onLoop() {
            uint16_t mv_val;
            int err = zbhome::sensors::read_adc_mv(&m_adc_spec, &mv_val);
            if (err < 0) {
                LOG_ERR("failed to fetch adc mv value: %d", err);
                return;
            }
        };
    }
}