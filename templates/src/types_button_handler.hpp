#pragma once

#include <vector>
#include <functional>

#include <zephyr/drivers/gpio.h>

#include "types.hpp"

namespace zbhome
{

    /*
    This header is always included in main, and ready to be used.

    To handle button0(or any gpio defined as button in dts) status changes you can use this header like this:
    ```
    uint32_t button_bit = BIT(DT_NODE_CHILD_IDX(DT_NODELABEL(button0)));
    addButtonHandler(button0, [](uint32_t button_bit, bool newStatus) {
        printf("new button state is %d", newStatus);
    });
    ```

    Handler function that is passed to `addButtonHandler`
    is pretty flexible thanks to c++ function type.
    */

    // Handler should not take long to complete, as handlers
    // will be run sequentially. So if there are 5 handlers
    // registered - they will be run one after another.
    typedef std::function<void(uint32_t button_bit, bool newState)> buttonHandlerFn;

    struct buttonHandlerConfiguration
    {
        uint32_t button_bit;
        buttonHandlerFn buttonHandler;
    };

    void button_changed(uint32_t button_state, uint32_t has_changed);

    bool setupButtonHandler(uint32_t factory_reset_button_bit, uint32_t identify_mode_button_bit);

    // Currently, it is possible to set up multiple handlers for the same button.
    void addButtonHandler(uint32_t button_bit, buttonHandlerFn handler);
}