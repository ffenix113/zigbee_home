#include <zephyr/logging/log.h>
#include <zephyr/drivers/gpio.h>

#include "clusters.hpp"
#ifdef __cplusplus
extern "C"
{
#endif
#include <zb_zcl_ias_zone.h>
#include <zigbee/zigbee_error_handler.h>
#ifdef __cplusplus
}
#endif

#include "types.hpp"
#include "types_ias_zone.hpp"
#include "types_button_handler.hpp"

#include <dk_buttons_and_leds.h>

// For some reason ZB_SCHEDULE_CALLBACK is not defined with current setup,
// and I can't find a necessary include/config to enable it.
// So for now - re-define the callback as another callback.
// #define ZB_SCHEDULE_CALLBACK ZB_SCHEDULE_APP_CALLBACK

// LOG_MODULE_DECLARE(app, LOG_LEVEL_INF);
LOG_MODULE_REGISTER(ias_zone, LOG_LEVEL_DBG);

namespace zbhome
{
    namespace components
    {
        IASZone::IASZone(uint32_t button_bit) : m_button_bit(button_bit) {
                                                };

        bool IASZone::setup()
        {
            // This code is not in constructor because endpoint is not yet set at that point.
            const uint8_t endpoint = getEndpoint();
            const auto handler = [endpoint](uint32_t button_bit, bool newStatus)
            {
                // A debug led to show that handler is triggered correctly.
                // dk_set_led(1, newStatus);

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

            addButtonHandler(m_button_bit, handler);

            return zbhome::types::Component::setup();
        }

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

            if (bufid)
            {
                zb_buf_free(bufid);
            }
        };
    }
}