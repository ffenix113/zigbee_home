#pragma once

#include <zephyr/logging/log.h>
#include <zephyr/drivers/gpio.h>

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
        class IASZone : public zbhome::types::Component
        {
            IASZone(const struct gpio_dt_spec pin);

        private:
            static void update_zone_status(zb_bufid_t bufid, zb_uint16_t cb_data);
        };
    }
}