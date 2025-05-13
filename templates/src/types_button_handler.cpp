#include <vector>
#include <functional>

#include <zephyr/drivers/gpio.h>

#include <zboss_api.h>

#include "types.hpp"
#include "types_button_handler.hpp"

namespace zbhome
{
    void buttonHandler(const struct device *port, struct gpio_callback *cb, gpio_port_pins_t pins)
    {
        // printf("button handler called. port %s, pins %d\n", port->name, pins);
        for (const auto &handlerCfg : attachedHandlers.handlers)
        {
            // printf("iter btn: %s/%d\n", handlerCfg.button.port->name, BIT(handlerCfg.button.pin));
            if (handlerCfg.button.port == port && (pins & BIT(handlerCfg.button.pin)) != 0)
            {
                handlerCfg.buttonHandler(&handlerCfg.button, gpio_pin_get_dt(&handlerCfg.button) != 0);
            }
        }
    }

    void setupButtonHandler()
    {
        gpio_init_callback(&attachedHandlers.cb, buttonHandler, -1);
    };

    void addButtonHandler(const struct gpio_dt_spec button, buttonHandlerFn handler)
    {
        attachedHandlers.handlers.push_back(buttonHandlerConfiguration{
            .button = button,
            .buttonHandler = handler,
        });

        gpio_add_callback_dt(&button, &attachedHandlers.cb);
    }
}