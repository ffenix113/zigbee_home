#pragma once

#include <zephyr/drivers/gpio.h>
#include <zephyr/input/input.h>
#include <zephyr/logging/log.h>

#include "types.hpp"

// For some reason ZB_SCHEDULE_CALLBACK is not defined with current setup,
// and I can't find a necessary include/config to enable it.
// So for now - re-define the callback as another callback.
#define ZB_SCHEDULE_CALLBACK ZB_SCHEDULE_APP_CALLBACK

// LOG_MODULE_DECLARE(app, LOG_LEVEL_INF);
// LOG_MODULE_REGISTER(ias_zone, LOG_LEVEL_INF);

namespace zbhome
{
    namespace components
    {
        class LevelControl : public zbhome::types::Component
        {
        public:
            // Currently this component is highly-specialized for rotary encoder use.
            // Feel free to make it more generic.
            LevelControl(uint8_t max_rotations_per_second);
            bool setup() override;
            // Only for input callback
            static void m_input_cb(input_event *evt, void *user_data);

        private:
            static void m_send_step_cmd(zb_bufid_t bufid, zb_uint16_t cb_data);

            // Default value to 2 reotations per second.
            const uint32_t m_rotation_time = 1000 / 2;
            uint32_t m_last_ms = 0;
            bool m_prev_direction = true;
        };
    }
}