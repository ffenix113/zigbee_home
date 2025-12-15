#pragma once

#include <zephyr/settings/settings.h>

namespace zbhome
{
    namespace settings
    {
        typedef void(settings_loaded_cb)();

        // Add setting to global list.
        // It will then be loaded and saved as necessary.
        //
        // This function must be called before `load` will be.
        //
        // Settings with duplicate names will be set, no checks for that.
        int8_t add_setting(const char *name, void *ptr, size_t size);
        // This will add callback for when all settings were loaded.
        // It is intended for cases when some code relies on settings be loaded,
        // but does not know that it will be executed after settings will be loaded.
        int8_t add_loaded_cb(settings_loaded_cb loaded_cb);

        int load();
        int save();

        // If necessary - functions can be added to save specific setting.

        bool get_enable_zigbee();
        // This option will return current value, while also setting it to the new_val.
        // This is useful to use the option only once.
        bool set_enable_zigbee(bool new_val);
    }
}