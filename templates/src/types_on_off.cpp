#include <zephyr/drivers/gpio.h>

#include <zboss_api.h>
#include <zb_zcl_common.h>

#include "types.hpp"
#include "types_zcl_command_handler.hpp"

namespace zbhome {
    namespace components {
        void OnOff::zclSetAttrValue(zb_zcl_set_attr_value_param_t* setAttrValueParam) {
            if (setAttrValueParam->cluster_id != ZB_ZCL_CLUSTER_ID_ON_OFF ||
                setAttrValueParam->attr_id != ZB_ZCL_ATTR_ON_OFF_ON_OFF_ID) {
                    return;
                }

            bool value = (bool)setAttrValueParam->values.data8;
            gpio_pin_set_dt(&m_pin, value);
        }
    }
}