
#include <zephyr/bluetooth/bluetooth.h>
#include <zephyr/bluetooth/hci.h>
#include <zephyr/bluetooth/conn.h>
#include <zephyr/bluetooth/uuid.h>
#include <zephyr/bluetooth/gatt.h>
#include <zephyr/dfu/mcuboot.h>
#include <zephyr/mgmt/mcumgr/mgmt/callbacks.h>
#include <zephyr/mgmt/mcumgr/grp/os_mgmt/os_mgmt_callbacks.h>

#include <zephyr/logging/log.h>
#include <zephyr/sys/reboot.h>

#include "extenders.hpp"
#include "ble_ota.hpp"
#include "settings.hpp"

LOG_MODULE_REGISTER(ble_ota, LOG_LEVEL_DBG);

#define OTABLE_DEVICE_NAME CONFIG_BT_DEVICE_NAME
#define OTABLE_DEVICE_NAME_LEN (sizeof(OTABLE_DEVICE_NAME) - 1)

#define OTABLE_ADV_MODIFIER 5
#define OTABLE_ADV_INT_MIN BT_GAP_ADV_FAST_INT_MIN_2 *OTABLE_ADV_MODIFIER
#define OTABLE_ADV_INT_MAX BT_GAP_ADV_FAST_INT_MAX_2 *OTABLE_ADV_MODIFIER

namespace zbhome
{
    namespace experimental
    {
        static const struct bt_data bleota_ad[] = {
            BT_DATA_BYTES(BT_DATA_FLAGS, (BT_LE_AD_GENERAL | BT_LE_AD_NO_BREDR)),
            BT_DATA(BT_DATA_NAME_COMPLETE, OTABLE_DEVICE_NAME, OTABLE_DEVICE_NAME_LEN),
        };

        static const struct bt_data bleota_sd[] = {};

        // Make advertising intervals longer, so we spend less power on it.
        const bt_le_adv_param bt_le_adv_params = BT_LE_ADV_PARAM_INIT(BT_LE_ADV_OPT_CONNECTABLE, // | BT_LE_ADV_OPT_ONE_TIME,
                                                                      OTABLE_ADV_INT_MIN, OTABLE_ADV_INT_MAX, NULL);

        static struct mgmt_callback bleota_dfu_cb;
        static struct mgmt_callback bleota_reset_cb;

        static void on_connect(struct bt_conn *conn, uint8_t err) {
            // LOG_INF("connected, changing security");
            // int err2 = bt_conn_set_security(conn, BT_SECURITY_L2);
            // if (err2 != 0)
            // {
            //     LOG_ERR("cannot set security: %d", err2);
            //     bt_conn_disconnect(conn, BT_HCI_ERR_AUTH_FAIL);
            // };
        };

        BT_CONN_CB_DEFINE(set_sec) = {
            .connected = on_connect,
        };

        int BLEOTA::on_main()
        {
            bleota_dfu_cb.callback = on_img_uploaded;
            bleota_dfu_cb.event_id = MGMT_EVT_OP_IMG_MGMT_DFU_PENDING;
            mgmt_callback_register(&bleota_dfu_cb);

            bleota_reset_cb.callback = on_reset_evt;
            bleota_reset_cb.event_id = MGMT_EVT_OP_OS_MGMT_RESET; // MGMT_EVT_OP_SETTINGS_MGMT_ACCESS;
            mgmt_callback_register(&bleota_reset_cb);

            if (int err = boot_write_img_confirmed(); err != 0)
            {
                LOG_ERR("confirm image: %d", err);
                return err;
            }

            int err = bt_enable(NULL);
            if (err)
            {
                LOG_ERR("Bluetooth init failed (err %d)", err);
                return err;
            }

            err = bt_le_adv_start(&bt_le_adv_params, bleota_ad, ARRAY_SIZE(bleota_ad),
                                  bleota_sd, ARRAY_SIZE(bleota_sd));
            if (err)
            {
                LOG_ERR("Advertising failed to start (err %d)", err);
                return err;
            }

            // load_without_zigbee = true;
            // if (int err = zbhome::settings::save(); err != 0)
            // {
            //     LOG_ERR("save settings: %d", err);
            // }

            return 0;
        };

        // When image is finished downloading - we want to mark image to be tested
        // and on first boot it would be set as permanent by BLEOTA::on_main.
        mgmt_cb_return BLEOTA::on_img_uploaded(uint32_t event, enum mgmt_cb_return prev_status,
                                               int32_t *rc, uint16_t *group, bool *abort_more, void *data,
                                               size_t data_size)
        {
            if (event != MGMT_EVT_OP_IMG_MGMT_DFU_PENDING)
            {
                LOG_INF("Chunk notification");
                return MGMT_CB_OK;
            }

            LOG_INF("DFU image fully uploaded");

            if (int err = boot_request_upgrade(BOOT_UPGRADE_TEST); err != 0)
            {
                LOG_ERR("request image upgrade: %d", err);
            }

            return MGMT_CB_OK;
        };

        mgmt_cb_return BLEOTA::on_reset_evt(uint32_t event, enum mgmt_cb_return prev_status,
                                            int32_t *rc, uint16_t *group, bool *abort_more, void *data,
                                            size_t data_size)
        {
            if (((struct os_mgmt_reset_data *)data)->force)
            {
                zbhome::settings::set_enable_zigbee(false);
                LOG_DBG("disabled zigbee on next boot via force reset");
                k_sleep(K_SECONDS(1));
            }

            return MGMT_CB_OK;
        };
    }
}

REGISTER_EXTENDER(bleota, zbhome::experimental::BLEOTA::on_main);
