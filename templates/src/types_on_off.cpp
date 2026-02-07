#include <zephyr/drivers/gpio.h>
#include <zephyr/logging/log.h>

#ifdef __cplusplus
extern "C"
{
#endif
#include <zboss_api.h>
#include <zcl/zb_zcl_on_off.h>
#include <zigbee/zigbee_error_handler.h>
#ifdef __cplusplus
}
#endif

#include "types.hpp"
#include "types_on_off.hpp"
#include "types_button_handler.hpp"

LOG_MODULE_REGISTER(on_off, LOG_LEVEL_DBG);

namespace zbhome
{
    namespace components
    {
        bool OnOff::setup()
        {
            if (m_button != 0)
            {
                const auto endpoint = getEndpoint();
                const auto output_pin = m_pin;

                const auto handler = [output_pin, endpoint](uint32_t button_bit, bool newState)
                {
                    // Only run this if button is pressed down, not released.
                    // Otherwise we will be tied to button state.
                    if (!newState)
                    {
                        return;
                    }

                    // Technically we could get this from `dev_ctx`,
                    // but we don't import `device.hpp`.
                    // Also, by calling this function we are not tied
                    // to any specific way the attributes are defined,
                    // or even if `dev_ctx` will be renamed.
                    const zb_zcl_attr_t *attr = zb_zcl_get_attr_desc_a(endpoint,
                                                                       ZB_ZCL_CLUSTER_ID_ON_OFF,
                                                                       ZB_ZCL_CLUSTER_SERVER_ROLE,
                                                                       ZB_ZCL_ATTR_ON_OFF_ON_OFF_ID);
                    bool current_state = *(bool *)attr->data_p;
                    // This is for readability only.
                    bool next_state = !current_state;

                    // This value should be next state, not current.
                    zb_uint16_t next_state_arg = next_state;
                    next_state_arg = (endpoint << 1) | next_state_arg;

                    gpio_pin_set_dt(&output_pin, next_state);

                    /* Allocate output buffer and set new state. */
                    zb_ret_t zb_err_code = zb_buf_get_out_delayed_ext(
                        OnOff::change_state, next_state_arg, 0);
                    ZB_ERROR_CHECK(zb_err_code);
                };

                zbhome::addButtonHandler(m_button, std::move(handler));
            }

            return zbhome::types::Component::setup();
        }

        void OnOff::zclSetAttrValue(zb_zcl_set_attr_value_param_t *setAttrValueParam)
        {
            if (setAttrValueParam->cluster_id != ZB_ZCL_CLUSTER_ID_ON_OFF ||
                setAttrValueParam->attr_id != ZB_ZCL_ATTR_ON_OFF_ON_OFF_ID)
            {
                return;
            }

            bool value = (bool)setAttrValueParam->values.data8;
            gpio_pin_set_dt(&m_pin, value);
        }

        void OnOff::change_state(zb_bufid_t bufid, zb_uint16_t data)
        {
            zb_uint8_t state_val = data & 1;
            const zb_uint8_t endpoint = data >> 1;

            auto status = zb_zcl_set_attr_val(endpoint,
                                              ZB_ZCL_CLUSTER_ID_ON_OFF,
                                              ZB_ZCL_CLUSTER_SERVER_ROLE,
                                              ZB_ZCL_ATTR_ON_OFF_ON_OFF_ID,
                                              &state_val,
                                              ZB_FALSE);
            if (status != 0)
            {
                LOG_DBG("set attr status: %d", status);
            }
        }
    }
}