#include <zephyr/logging/log.h>
#include <zephyr/drivers/gpio.h>

#include <zboss_api.h>
#include <zb_zcl_ias_zone.h>

#include <zigbee/zigbee_error_handler.h>

#include "clusters.hpp"
#include "types.hpp"
#include "types_ias_zone.hpp"
#include "types_button_handler.hpp"

// For some reason ZB_SCHEDULE_CALLBACK is not defined with current setup,
// and I can't find a necessary include/config to enable it.
// So for now - re-define the callback as another callback.
#define ZB_SCHEDULE_CALLBACK ZB_SCHEDULE_APP_CALLBACK

// LOG_MODULE_DECLARE(app, LOG_LEVEL_INF);
LOG_MODULE_REGISTER(ias_zone, LOG_LEVEL_INF);

namespace zbhome
{
    namespace components
    {
        IASZone::IASZone(const struct gpio_dt_spec pin)
        {
            gpio_pin_configure_dt(&pin, GPIO_INPUT);
            gpio_pin_interrupt_configure_dt(&pin, GPIO_INT_EDGE_BOTH);

            const uint8_t endpoint = getEndpoint();
            const auto handler = [endpoint](const struct gpio_dt_spec *btn, bool newStatus)
            {
                // Pack data.
                // Endpoint can be >127, so we can't pack it and state into single uint8.
                //
                // First bit is the status of the button,
                // next 8 bits - endpoint number.
                zb_uint16_t data = endpoint << 1 | newStatus;

                /* Allocate output buffer and send on/off command. */
                zb_ret_t zb_err_code = zb_buf_get_out_delayed_ext(
                    IASZone::update_zone_status, data, 0);
                ZB_ERROR_CHECK(zb_err_code);
            };

            addButtonHandler(pin, handler);
        };

        void IASZone::update_zone_status(zb_bufid_t bufid, zb_uint16_t cb_data)
        {
            // TODO: Probably this function needs to free bufid somewhere,
            // but I am not sure if it is actually the case.
            // Would need to double-check.

            // Decode values from the argument.
            bool status = cb_data & 1;
            zb_uint8_t endpoint = cb_data >> 1;

            switch (status)
            {
            case true:
                ZB_ZCL_IAS_ZONE_SET_BITS(bufid, endpoint, ZB_ZCL_IAS_ZONE_ZONE_STATUS_ALARM1);
                break;
            case false:
                ZB_ZCL_IAS_ZONE_CLEAR_BITS(bufid, endpoint, ZB_ZCL_IAS_ZONE_ZONE_STATUS_ALARM1);
                break;
            }
        };
    }
}