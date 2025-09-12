#pragma once

#include <zephyr/drivers/adc.h>

#include "types.hpp"

namespace zbhome {
    namespace components {
        class PowerConfig: public zbhome::types::Sensor {
            public:
                PowerConfig(const struct adc_dt_spec adc_spec, uint16_t min_mv_threshold, uint16_t rated_mv): 
                    m_adc_spec(adc_spec), 
                    m_min_mv_threshold(min_mv_threshold),
                    m_rated_mv(rated_mv)
                {};
                void onLoop();

            private:
                const struct adc_dt_spec m_adc_spec = {};

                uint16_t m_min_mv_threshold = 0;
                uint16_t m_rated_mv = 0;
        };
    }
}