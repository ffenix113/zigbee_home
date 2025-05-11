#pragma once

#include <vector>
#include <functional>

#include <zephyr/drivers/gpio.h>

#include "types.hpp"


/*
This header is always included in main, and ready to be used.

To handle button(gpio input) changes you can use this header like this:
```
struct gpio_dt_spec button0 = GPIO_DT_SPEC_GET(DT_NODELABEL(button0), gpios);
gpio_pin_configure_dt(&button0, GPIO_INPUT);
gpio_pin_interrupt_configure_dt(&button0, GPIO_INT_EDGE_BOTH);

addButtonHandler(button0, [](const struct gpio_dt_spec* btn, bool newStatus) {
    printf("new button state is %d", newStatus);
});
```

Handler function that is passed to `addButtonHandler` 
is pretty flexible thanks to c++ function type.
*/


// Handler should not take long to complete, as handlers
// will be run sequentially. So if there are 5 handlers
// registered - they will be run one after another.
typedef std::function<void(const struct gpio_dt_spec *btn, bool newState)> buttonHandlerFn;

struct buttonHandlerConfiguration{
    const struct gpio_dt_spec button;
    buttonHandlerFn buttonHandler;
};

struct {
    struct gpio_callback cb;
    std::vector<struct buttonHandlerConfiguration> handlers;
} attachedHandlers;


void buttonHandler(const struct device *port, struct gpio_callback *cb, gpio_port_pins_t pins) {
    // printf("button handler called. port %s, pins %d\n", port->name, pins);
    for (const auto& handlerCfg : attachedHandlers.handlers) {
        // printf("iter btn: %s/%d\n", handlerCfg.button.port->name, BIT(handlerCfg.button.pin));
        if (handlerCfg.button.port == port && (pins & BIT(handlerCfg.button.pin)) != 0) {
            handlerCfg.buttonHandler(&handlerCfg.button, gpio_pin_get_dt(&handlerCfg.button) != 0);
        }
    }
}

void setupButtonHandler() {
    gpio_init_callback(&attachedHandlers.cb, buttonHandler, -1);
};

void addButtonHandler(const struct gpio_dt_spec button, buttonHandlerFn handler) {
    attachedHandlers.handlers.push_back(buttonHandlerConfiguration{
        .button = button,
        .buttonHandler = handler,
    });

    gpio_add_callback_dt(&button, &attachedHandlers.cb);
}
