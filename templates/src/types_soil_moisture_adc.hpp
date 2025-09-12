#pragma once

#include <zephyr/drivers/adc.h>

#include "types.hpp"

namespace zbhome
{
    namespace components
    {
        class SoilMoistureADC : public zbhome::types::Sensor
        {
        public:
            SoilMoistureADC(const struct adc_dt_spec adc_spec, uint16_t min_mv_val, uint16_t max_mv_val) : m_adc_spec(adc_spec),
                                                                                                           m_min_mv_val(min_mv_val),
                                                                                                           m_max_mv_val(max_mv_val) {};
            void onLoop();

        private:
            const struct adc_dt_spec m_adc_spec = {};

            uint16_t m_min_mv_val = 0;
            uint16_t m_max_mv_val = -1; // -1 will be set to max uint anyway
        };
    }
}