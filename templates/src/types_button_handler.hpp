#pragma once

#include <vector>
#include <functional>

#include <zephyr/drivers/gpio.h>

#include "types.hpp"

namespace zbhome
{

    /*
    This header is always included in main, and ready to be used.

    To handle button(gpio input) status changes you can use this header like this:
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

    struct buttonHandlerConfiguration
    {
        const struct gpio_dt_spec button;
        buttonHandlerFn buttonHandler;
    };

    struct
    {
        struct gpio_callback cb;
        std::vector<struct buttonHandlerConfiguration> handlers;
    } attachedHandlers;

    void buttonHandler(const struct device *port, struct gpio_callback *cb, gpio_port_pins_t pins);

    void setupButtonHandler();
    
    // Currently, it is possible to set up multiple handlers for the same gpio.
    void addButtonHandler(const struct gpio_dt_spec button, buttonHandlerFn handler);
}