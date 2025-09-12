#pragma once

#include <zephyr/drivers/gpio.h>

#include <zb_zcl_common.h>

#include "types.hpp"

namespace zbhome
{
    namespace components
    {
        class OnOff : public zbhome::types::ZCLCommandHandler
        {
        public:
            OnOff(const struct gpio_dt_spec pin) : m_pin(pin);
            // This function will be called when callback is received with the endpoint
            // that this components is for. Otherwise cluster/attribute id's are not
            // checked and are up to the implementation to validate.
            void zclSetAttrValue(zb_zcl_set_attr_value_param_t* setAttrValueParam);

        private:
            const struct gpio_dt_spec m_pin = {};
        };
    }
}