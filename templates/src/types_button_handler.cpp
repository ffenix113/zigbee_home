#include <vector>
#include <functional>

#include <zephyr/logging/log.h>

#include <dk_buttons_and_leds.h>
#include <zigbee/zigbee_app_utils.h>

#include "types_button_handler.hpp"

LOG_MODULE_REGISTER(button_handler, LOG_LEVEL_DBG);

namespace zbhome
{
    std::vector<struct buttonHandlerConfiguration> attached_handlers;
    uint32_t m_identify_mode_button_bit = 0;

    // Having this setup, as shown by example Nordic Zigbee devide simplifies
    // handling a bit, and saves some memory as well.
    //
    // Additional benefit that it is by-default "wired" for zigbee factory reset
    // and other helpful things.
    void button_changed(uint32_t button_state, uint32_t has_changed)
    {
        LOG_DBG("button_changed: state: 0x%08u, changed: 0x%08u", button_state, has_changed);

#if CONFIG_ZIGBEE_ROLE_END_DEVICE
        /* Inform default signal handler about user input at the device */
        user_input_indicate();
#endif

        if (m_identify_mode_button_bit & has_changed)
        {
            if (m_identify_mode_button_bit & button_state)
            {
                /* Button changed its state to pressed */
            }
            else
            {
                /* Button changed its state to released */
                if (was_factory_reset_done())
                {
                    /* The long press was for Factory Reset */
                    // LOG_DBG("After Factory Reset - ignore button release");
                }
                else
                {
                    /* Button released before Factory Reset */

                    /* Start identification mode */
                    // ZB_SCHEDULE_APP_CALLBACK(start_identifying, 0);
                }
            }
        }

        for (const auto &handlerCfg : attached_handlers)
        {
            if (handlerCfg.button_bit & has_changed)
            {
                handlerCfg.buttonHandler(handlerCfg.button_bit, handlerCfg.button_bit & button_state);
            }
        }

        check_factory_reset_button(button_state, has_changed);
    }

    bool setupButtonHandler(uint32_t factory_reset_button_bit, uint32_t identify_mode_button_bit)
    {

        int err = 0;
        err = dk_buttons_init(button_changed);
        if (err)
        {
            LOG_ERR("Cannot init buttons (err: %d)", err);
            return false;
        }

        register_factory_reset_button(factory_reset_button_bit);
        m_identify_mode_button_bit = identify_mode_button_bit;

        return true;
    };

    void addButtonHandler(uint32_t button_bit, buttonHandlerFn handler)
    {
        attached_handlers.push_back(buttonHandlerConfiguration{
            .button_bit = button_bit,
            .buttonHandler = handler,
        });
    }
}