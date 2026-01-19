#include <zephyr/logging/log.h>
#include <zephyr/input/input.h>

#include "clusters.hpp"
#ifdef __cplusplus
extern "C"
{
#endif
#include <zb_zcl_level_control.h>
#include <zigbee/zigbee_error_handler.h>
#ifdef __cplusplus
}
#endif

#include <dk_buttons_and_leds.h>

#include "types.hpp"
#include "types_level_control.hpp"
#include "types_button_handler.hpp"

LOG_MODULE_REGISTER(level, LOG_LEVEL_DBG);

namespace zbhome
{
    namespace components
    {
        // A bit of a hack, because we need to match device to levelControl instance.
        static std::unordered_map<const struct device *, LevelControl *> level_map = {};
        static const auto level_input_cb = LevelControl::m_input_cb;
        // This callback will be defined as static block in specific ROM location,
        // so we can't pass in runtime arguments.
        INPUT_CALLBACK_DEFINE(nullptr, level_input_cb, nullptr);

        LevelControl::LevelControl(uint8_t max_rotations_per_second) : m_rotation_time(1000 / max_rotations_per_second) {
                                                                       };

        bool LevelControl::setup()
        {
            // This code is not in constructor because device is not yet set at that point.
            level_map[this->getDevice()] = this;

            return zbhome::types::Component::setup();
        }

        void LevelControl::m_input_cb(input_event *evt, void *user_data)
        {
            if (evt->code != INPUT_REL_WHEEL)
            {
                return;
            }

            LevelControl *level_inst = level_map[evt->dev];
            // Check if we want to listen to this event.
            if (level_inst == nullptr)
            {
                return;
            }

            const bool same_direction = (evt->value >= 1) == (bool)level_inst->m_prev_direction;
            // const uint32_t uptime = k_cycle_get_32();
            const uint32_t uptime = (uint32_t)k_uptime_get();

            const uint32_t time_diff = uptime - level_inst->m_last_ms;
            // LOG_DBG("%u, %u, %u", uptime, uptime - level_inst->m_last_ms, level_inst->m_rotation_time);
            if (same_direction && (time_diff < level_inst->m_rotation_time))
            {
                // LOG_DBG("level change skipped %p", level_inst);
                return;
            }

            const zb_uint8_t endpoint = level_inst->getEndpoint();
            LOG_DBG("level change %u: %d", time_diff, evt->value);
            const zb_uint16_t data = endpoint << 1 | evt->value >= 1;

            level_inst->m_last_ms = uptime;
            level_inst->m_prev_direction = evt->value >= 1;

            /* Allocate output buffer and send on/off command. */
            zb_ret_t zb_err_code = zb_buf_get_out_delayed_ext(m_send_step_cmd, data, 0);
            ZB_ERROR_CHECK(zb_err_code);
        }

        void LevelControl::m_send_step_cmd(zb_bufid_t bufid, zb_uint16_t cb_data)
        {
            // TODO: Probably this function needs to free bufid somewhere,
            // but I am not sure if it is actually the case.
            // Would need to double-check.

            // Decode values from the argument.
            bool direction = cb_data & 1;
            zb_uint8_t endpoint = cb_data >> 1;

            // zb_zcl_level_control_step_req_t *req;

            const zb_uint8_t addr_ep = ZB_APS_ADDR_MODE_DST_ADDR_ENDP_NOT_PRESENT;

            zb_zcl_level_control_step_mode_e step_mode = direction ? ZB_ZCL_LEVEL_CONTROL_STEP_MODE_UP : ZB_ZCL_LEVEL_CONTROL_STEP_MODE_DOWN;

            ZB_ZCL_LEVEL_CONTROL_SEND_STEP_REQ_ZCL8(bufid, addr_ep, 0, 0,
                                                    endpoint, ZB_AF_HA_PROFILE_ID, ZB_ZCL_DISABLE_DEFAULT_RESPONSE,
                                                    NULL, step_mode,
                                                    1, 0, 0, 0);

            // Buffer here should not be freed.
            // Otherwise it crashes ZBOSS stack (for some reason).
        };
    }
}