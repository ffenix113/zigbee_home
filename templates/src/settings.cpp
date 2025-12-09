#include <stdio.h>

#include <zephyr/logging/log.h>
#include <zephyr/settings/settings.h>

#include "settings.hpp"

LOG_MODULE_REGISTER(zbhome_settings, LOG_LEVEL_DBG);

namespace zbhome
{
    namespace settings
    {
        typedef struct
        {
            const char *name;
            void *ptr;
            size_t size;
        } setting;

        // our own settings
        static bool p_should_enable_zigbee = true;

        // Just any value to not have too much settings here.
        // It can be increased later.
        constexpr uint8_t max_settings = 4;
        // Number of default settings present in the array bellow.
        static uint8_t settings_count = 1;
        static setting settings[max_settings] = {
            {
                .name = "enable_zigbee",
                .ptr = &p_should_enable_zigbee,
                .size = sizeof(p_should_enable_zigbee),
            }};

        static uint8_t settings_loaded_cb_count = 0;
        constexpr uint8_t max_cb = max_settings * 2;
        static settings_loaded_cb *loaded_cbs[max_cb];

        // Settings with duplicate names will be set, no checks for that.
        // This function returns index of the inserted setting.
        //
        // If return value is negative - setting is not added to the list.
        int8_t add_setting(const char *name, void *ptr, size_t size)
        {
            if (settings_count == max_settings)
            {
                return -1;
            }

            settings[settings_count] = {
                .name = name,
                .ptr = ptr,
                .size = size,
            };

            settings_count++;

            return settings_count - 1;
        }

        int8_t add_loaded_cb(settings_loaded_cb cb)
        {
            if (settings_loaded_cb_count == max_cb)
            {
                return -1;
            }

            loaded_cbs[settings_loaded_cb_count] = cb;
            settings_loaded_cb_count++;

            return settings_loaded_cb_count - 1;
        }

        int settings_handler_set(const char *name, size_t len, settings_read_cb read_cb, void *cb_arg)
        {
            const char *next;
            int err;

            for (int i = 0; i < settings_count; i++)
            {
                auto const &setting = settings[i];

                if (settings_name_steq(name, setting.name, &next) && !next)
                {
                    if (len != setting.size)
                    {
                        LOG_ERR("setting %s size to set is invalid: want %d, got %d", name, setting.size, len);
                        return -EINVAL;
                    }

                    err = read_cb(cb_arg, setting.ptr, setting.size);
                    if (err <= 0)
                    {
                        return err;
                    }

                    LOG_DBG("loaded %s", name);
                }
            }

            LOG_INF("%d setting(s) loaded", settings_count);

            return 0;
        }

        int settings_handler_export(int(export_fn)(const char *name, const void *val,
                                                   size_t val_len))
        {
            for (int i = 0; i < settings_count; i++)
            {
                auto const &setting = settings[i];

                if (int err = export_fn(setting.name, setting.ptr, setting.size); err != 0)
                {
                    LOG_ERR("save %s: %d", setting.name, err);
                    return err;
                }
            }

            return 0;
        }

        int settings_handler_commit(void)
        {
            for (int i = 0; i < settings_loaded_cb_count; i++)
            {
                auto const &cb = loaded_cbs[i];
                cb();
            }

            return 0;
        }

        static struct settings_handler zbhome_conf = {
            .name = "zbhome",
            .h_set = settings_handler_set,
            .h_commit = settings_handler_commit,
            .h_export = settings_handler_export,
        };

        int load()
        {
            int err = settings_subsys_init();
            if (err)
            {
                LOG_ERR("settings_subsys_init failed (err %d)", err);
            }

            err = settings_register(&zbhome_conf);
            if (err)
            {
                LOG_ERR("settings_register failed (err %d)", err);
            }

            err = settings_load_subtree(zbhome_conf.name);
            if (err)
            {
                LOG_ERR("settings_load failed (err %d)", err);
            }

            return 0;
        }

        int save()
        {
            constexpr uint8_t buf_size = 64;
            char buf[buf_size];

            for (int i = 0; i < settings_count; i++)
            {
                auto const &setting = settings[i];

                snprintf(buf, buf_size, "zbhome/%s", setting.name);
                if (int err = settings_save_one(buf, setting.ptr, setting.size); err != 0)
                {
                    LOG_ERR("save %s: %d", buf, err);
                    return err;
                }
            }

            return 0;
        }

        bool get_enable_zigbee()
        {
            return p_should_enable_zigbee;
        }

        bool set_enable_zigbee(bool new_val)
        {
            bool current = get_enable_zigbee();
            p_should_enable_zigbee = new_val;

            if (new_val != current)
            {
                save();
            }

            return current;
        }
    }
}
