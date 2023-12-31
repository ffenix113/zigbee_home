/*
 * Copyright (c) 2022 Nordic Semiconductor ASA
 *
 * SPDX-License-Identifier: LicenseRef-Nordic-5-Clause
 */

#include <zephyr/device.h>
#include <dk_buttons_and_leds.h>
#include <zephyr/drivers/uart.h>
#include <zephyr/logging/log.h>
#include <ram_pwrdn.h>
#include <zb_nrf_platform.h>
#include <zboss_api.h>
#include <zboss_api_addons.h>
#include <zephyr/kernel.h>
#include <zigbee/zigbee_app_utils.h>
#include <zigbee/zigbee_error_handler.h>

#ifdef CONFIG_USB_DEVICE_STACK
#include <zephyr/usb/usb_device.h>
#include <zephyr/usb/usbd.h>
#include <zephyr/sys/printk.h>
#endif /* CONFIG_USB_DEVICE_STACK */

#include "weather_station.h"

/* Delay for console initialization */
#define WAIT_FOR_CONSOLE_MSEC 100
#define WAIT_FOR_CONSOLE_DEADLINE_MSEC 500

/* Weather check period */
#define WEATHER_CHECK_PERIOD_MSEC (1000 * CONFIG_WEATHER_CHECK_PERIOD_SECONDS)

/* Delay for first weather check */
#define WEATHER_CHECK_INITIAL_DELAY_MSEC (1000 * CONFIG_FIRST_WEATHER_CHECK_DELAY_SECONDS)

/* Time of LED on state while blinking for identify mode */
#define IDENTIFY_LED_BLINK_TIME_MSEC 500

#define LED_POWER DK_LED1
/* In Thingy53 each LED is a RGB component of a single LED */
#define LED_RED DK_LED2
#define LED_GREEN DK_LED3
#define LED_BLUE DK_LED4

/* LED indicating that device successfully joined Zigbee network */
#define ZIGBEE_NETWORK_STATE_LED LED_BLUE

/* LED used for device identification */
#define IDENTIFY_LED LED_RED

/* Button used to enter the Identify mode */
#define IDENTIFY_MODE_BUTTON DK_BTN1_MSK

/* Button to start Factory Reset */
#define FACTORY_RESET_BUTTON IDENTIFY_MODE_BUTTON

BUILD_ASSERT(DT_NODE_HAS_COMPAT(DT_CHOSEN(zephyr_console),
				zephyr_cdc_acm_uart),
	     "Console device is not ACM CDC UART device");
LOG_MODULE_REGISTER(app, LOG_LEVEL_DBG);// CONFIG_ZIGBEE_WEATHER_STATION_LOG_LEVEL);

/* Stores all cluster-related attributes */
static struct zb_device_ctx dev_ctx;

/* Attributes setup */
ZB_ZCL_DECLARE_BASIC_ATTRIB_LIST_EXT(
	basic_attr_list,
	&dev_ctx.basic_attr.zcl_version,
	&dev_ctx.basic_attr.app_version,
	&dev_ctx.basic_attr.stack_version,
	&dev_ctx.basic_attr.hw_version,
	dev_ctx.basic_attr.mf_name,
	dev_ctx.basic_attr.model_id,
	dev_ctx.basic_attr.date_code,
	&dev_ctx.basic_attr.power_source,
	dev_ctx.basic_attr.location_id,
	&dev_ctx.basic_attr.ph_env,
	dev_ctx.basic_attr.sw_ver);

/* Declare attribute list for Identify cluster (client). */
ZB_ZCL_DECLARE_IDENTIFY_CLIENT_ATTRIB_LIST(
	identify_client_attr_list);

/* Declare attribute list for Identify cluster (server). */
ZB_ZCL_DECLARE_IDENTIFY_SERVER_ATTRIB_LIST(
	identify_server_attr_list,
	&dev_ctx.identify_attr.identify_time);

ZB_ZCL_DECLARE_TEMP_MEASUREMENT_ATTRIB_LIST(
	temperature_measurement_attr_list,
	&dev_ctx.temp_attrs.measure_value,
	&dev_ctx.temp_attrs.min_measure_value,
	&dev_ctx.temp_attrs.max_measure_value,
	&dev_ctx.temp_attrs.tolerance
	);

ZB_ZCL_DECLARE_PRESSURE_MEASUREMENT_ATTRIB_LIST(
	pressure_measurement_attr_list,
	&dev_ctx.pres_attrs.measure_value,
	&dev_ctx.pres_attrs.min_measure_value,
	&dev_ctx.pres_attrs.max_measure_value,
	&dev_ctx.pres_attrs.tolerance
	);

ZB_ZCL_DECLARE_REL_HUMIDITY_MEASUREMENT_ATTRIB_LIST(
	humidity_measurement_attr_list,
	&dev_ctx.humidity_attrs.measure_value,
	&dev_ctx.humidity_attrs.min_measure_value,
	&dev_ctx.humidity_attrs.max_measure_value
	);

/* On/Off cluster attributes additions data */
ZB_ZCL_DECLARE_ON_OFF_ATTRIB_LIST(
	on_off_attr_list,
	&dev_ctx.on_off_attr.on_off);

/* Clusters setup */
ZB_HA_DECLARE_WEATHER_STATION_CLUSTER_LIST(
	weather_station_cluster_list,
	basic_attr_list,
	identify_client_attr_list,
	identify_server_attr_list,
	temperature_measurement_attr_list,
	pressure_measurement_attr_list,
	humidity_measurement_attr_list,
	on_off_attr_list);

/* Endpoint setup (single) */
ZB_HA_DECLARE_WEATHER_STATION_EP(
	weather_station_ep,
	WEATHER_STATION_ENDPOINT_NB,
	weather_station_cluster_list);

/* Device context */
ZBOSS_DECLARE_DEVICE_CTX_1_EP(
	weather_station_ctx,
	weather_station_ep);


/* Manufacturer name (32 bytes). */
#define BULB_INIT_BASIC_MANUF_NAME      "FFenix113"

/* Model number assigned by manufacturer (32-bytes long string). */
#define BULB_INIT_BASIC_MODEL_ID        "dongle"

/* First 8 bytes specify the date of manufacturer of the device
 * in ISO 8601 format (YYYYMMDD). The rest (8 bytes) are manufacturer specific.
 */
#define BULB_INIT_BASIC_DATE_CODE       "20200329"

/* Describes the physical location of the device (16 bytes).
 * May be modified during commissioning process.
 */
#define BULB_INIT_BASIC_LOCATION_DESC   ""
/* Describes the type of physical environment.
 * For possible values see section 3.2.2.2.10 of ZCL specification.
 */
#define BULB_INIT_BASIC_PH_ENV          ZB_ZCL_BASIC_ENV_UNSPECIFIED

static void mandatory_clusters_attr_init(void)
{
	/* Basic cluster attributes */
	dev_ctx.basic_attr.zcl_version = ZB_ZCL_VERSION;
	dev_ctx.basic_attr.power_source = ZB_ZCL_BASIC_POWER_SOURCE_DC_SOURCE;
	// Extended attributes
	dev_ctx.basic_attr.app_version = 01;
	dev_ctx.basic_attr.stack_version = 10;
	dev_ctx.basic_attr.hw_version = 11;

	/* Use ZB_ZCL_SET_STRING_VAL to set strings, because the first byte
	 * should contain string length without trailing zero.
	 *
	 * For example "test" string will be encoded as:
	 *   [(0x4), 't', 'e', 's', 't']
	 */
	ZB_ZCL_SET_STRING_VAL(
		dev_ctx.basic_attr.mf_name,
		BULB_INIT_BASIC_MANUF_NAME,
		ZB_ZCL_STRING_CONST_SIZE(BULB_INIT_BASIC_MANUF_NAME));

	ZB_ZCL_SET_STRING_VAL(
		dev_ctx.basic_attr.model_id,
		BULB_INIT_BASIC_MODEL_ID,
		ZB_ZCL_STRING_CONST_SIZE(BULB_INIT_BASIC_MODEL_ID));

	ZB_ZCL_SET_STRING_VAL(
		dev_ctx.basic_attr.date_code,
		BULB_INIT_BASIC_DATE_CODE,
		ZB_ZCL_STRING_CONST_SIZE(BULB_INIT_BASIC_DATE_CODE));

	ZB_ZCL_SET_STRING_VAL(
		dev_ctx.basic_attr.location_id,
		BULB_INIT_BASIC_LOCATION_DESC,
		ZB_ZCL_STRING_CONST_SIZE(BULB_INIT_BASIC_LOCATION_DESC));

	dev_ctx.basic_attr.ph_env = BULB_INIT_BASIC_PH_ENV;

	/* Identify cluster attributes data. */
	dev_ctx.identify_attr.identify_time =
		ZB_ZCL_IDENTIFY_IDENTIFY_TIME_DEFAULT_VALUE;

	/* Identify cluster attributes */
	dev_ctx.identify_attr.identify_time = ZB_ZCL_IDENTIFY_IDENTIFY_TIME_DEFAULT_VALUE;
}

static void measurements_clusters_attr_init(void)
{
	/* Temperature */
	dev_ctx.temp_attrs.measure_value = ZB_ZCL_ATTR_TEMP_MEASUREMENT_VALUE_UNKNOWN;
	dev_ctx.temp_attrs.min_measure_value = WEATHER_STATION_ATTR_TEMP_MIN;
	dev_ctx.temp_attrs.max_measure_value = WEATHER_STATION_ATTR_TEMP_MAX;
	dev_ctx.temp_attrs.tolerance = WEATHER_STATION_ATTR_TEMP_TOLERANCE;

	/* Pressure */
	dev_ctx.pres_attrs.measure_value = ZB_ZCL_ATTR_PRESSURE_MEASUREMENT_VALUE_UNKNOWN;
	dev_ctx.pres_attrs.min_measure_value = WEATHER_STATION_ATTR_PRESSURE_MIN;
	dev_ctx.pres_attrs.max_measure_value = WEATHER_STATION_ATTR_PRESSURE_MAX;
	dev_ctx.pres_attrs.tolerance = WEATHER_STATION_ATTR_PRESSURE_TOLERANCE;

	/* Humidity */
	dev_ctx.humidity_attrs.measure_value = ZB_ZCL_ATTR_REL_HUMIDITY_MEASUREMENT_VALUE_UNKNOWN;
	dev_ctx.humidity_attrs.min_measure_value = WEATHER_STATION_ATTR_HUMIDITY_MIN;
	dev_ctx.humidity_attrs.max_measure_value = WEATHER_STATION_ATTR_HUMIDITY_MAX;
	/* Humidity measurements tolerance is not supported at the moment */
}

static void toggle_identify_led(zb_bufid_t bufid)
{
	static bool led_on;

	led_on = !led_on;
	dk_set_led(IDENTIFY_LED, led_on);
	zb_ret_t err = ZB_SCHEDULE_APP_ALARM(toggle_identify_led,
					     bufid,
					     ZB_MILLISECONDS_TO_BEACON_INTERVAL(
						     IDENTIFY_LED_BLINK_TIME_MSEC));
	if (err) {
		LOG_ERR("Failed to schedule app alarm: %d", err);
	}
}

static void start_identifying(zb_bufid_t bufid)
{
	ZVUNUSED(bufid);

	if (ZB_JOINED()) {
		/*
		 * Check if endpoint is in identifying mode,
		 * if not put desired endpoint in identifying mode.
		 */
		if (dev_ctx.identify_attr.identify_time ==
		    ZB_ZCL_IDENTIFY_IDENTIFY_TIME_DEFAULT_VALUE) {

			zb_ret_t zb_err_code = zb_bdb_finding_binding_target(
				WEATHER_STATION_ENDPOINT_NB);

			if (zb_err_code == RET_OK) {
				dk_set_led_off(LED_POWER);
				dev_ctx.on_off_attr.on_off = 1;

				LOG_INF("Manually enter identify mode");
			} else if (zb_err_code == RET_INVALID_STATE) {
				LOG_WRN("RET_INVALID_STATE - Cannot enter identify mode");
			} else {
				ZB_ERROR_CHECK(zb_err_code);
			}
		} else {
			dk_set_led_on(LED_POWER);
			dev_ctx.on_off_attr.on_off = 0;
			LOG_INF("Manually cancel identify mode");
			zb_bdb_finding_binding_target_cancel();
		}
	} else {
		LOG_WRN("Device not in a network - cannot identify itself");
	}
}

static void identify_callback(zb_bufid_t bufid)
{
	zb_ret_t err = RET_OK;

	if (bufid) {
		/* Schedule a self-scheduling function that will toggle the LED */
		err = ZB_SCHEDULE_APP_CALLBACK(toggle_identify_led, bufid);
		if (err) {
			LOG_ERR("Failed to schedule app callback: %d", err);
		} else {
			LOG_INF("Enter identify mode");
		}
	} else {
		/* Cancel the toggling function alarm and turn off LED */
		err = ZB_SCHEDULE_APP_ALARM_CANCEL(toggle_identify_led,
						   ZB_ALARM_ANY_PARAM);
		if (err) {
			LOG_ERR("Failed to schedule app alarm cancel: %d", err);
		} else {
			dk_set_led_off(IDENTIFY_LED);
			LOG_INF("Cancel identify mode");
		}
	}
}

static void button_changed(uint32_t button_state, uint32_t has_changed)
{
	if (IDENTIFY_MODE_BUTTON & has_changed) {
		if (IDENTIFY_MODE_BUTTON & button_state) {
			/* Button changed its state to pressed */
		} else {
			/* Button changed its state to released */
			if (was_factory_reset_done()) {
				/* The long press was for Factory Reset */
				LOG_DBG("After Factory Reset - ignore button release");
			} else   {
				/* Button released before Factory Reset */

				/* Start identification mode */
				zb_ret_t err = ZB_SCHEDULE_APP_CALLBACK(start_identifying, 0);

				if (err) {
					LOG_ERR("Failed to schedule app callback: %d", err);
				}

				/* Inform default signal handler about user input at the device */
				user_input_indicate();
			}
		}
	}

	check_factory_reset_button(button_state, has_changed);
}

static void gpio_init(void)
{
	int err = dk_buttons_init(button_changed);

	if (err) {
		LOG_ERR("Cannot init buttons (err: %d)", err);
	}

	err = dk_leds_init();
	if (err) {
		LOG_ERR("Cannot init LEDs (err: %d)", err);
	}
}

#ifdef CONFIG_USB_DEVICE_STACK
static void wait_for_console(void)
{
	const struct device *console = DEVICE_DT_GET(DT_CHOSEN(zephyr_console));
	uint32_t dtr = 0;
	uint32_t time = 0;

	/* Enable the USB subsystem and associated HW */
	if (usb_enable(NULL)) {
		LOG_ERR("Failed to enable USB");
	} else {
		/* Wait for DTR flag or deadline (e.g. when USB is not connected) */
		while (!dtr && time < WAIT_FOR_CONSOLE_DEADLINE_MSEC) {
			uart_line_ctrl_get(console, UART_LINE_CTRL_DTR, &dtr);
			/* Give CPU resources to low priority threads */
			k_sleep(K_MSEC(WAIT_FOR_CONSOLE_MSEC));
			time += WAIT_FOR_CONSOLE_MSEC;
		}
		// while (!dtr) {
		// 	uart_line_ctrl_get(console, UART_LINE_CTRL_DTR, &dtr);
		// 	/* Give CPU resources to low priority threads. */
		// 	k_sleep(K_MSEC(500));
		// }
	}
}
#endif /* CONFIG_USB_DEVICE_STACK */

static void check_weather(zb_bufid_t bufid)
{
	ZVUNUSED(bufid);

	int err = weather_station_check_weather();

	if (err) {
		LOG_ERR("Failed to check weather: %d", err);
	} else {
		err = weather_station_update_temperature();
		if (err) {
			LOG_ERR("Failed to update temperature: %d", err);
		}

		err = weather_station_update_pressure();
		if (err) {
			LOG_ERR("Failed to update pressure: %d", err);
		}

		err = weather_station_update_humidity();
		if (err) {
			LOG_ERR("Failed to update humidity: %d", err);
		}
	}

	zb_ret_t zb_err = ZB_SCHEDULE_APP_ALARM(check_weather,
						0,
						ZB_MILLISECONDS_TO_BEACON_INTERVAL(
							WEATHER_CHECK_PERIOD_MSEC));
	if (zb_err) {
		LOG_ERR("Failed to schedule app alarm: %d", zb_err);
	}
}

/**@brief Callback function for handling ZCL commands.
 *
 * @param[in]   bufid   Reference to Zigbee stack buffer
 *                      used to pass received data.
 */
static void zcl_device_cb(zb_bufid_t bufid)
{
	zb_uint8_t cluster_id;
	zb_uint8_t attr_id;
	zb_zcl_device_callback_param_t  *device_cb_param =
		ZB_BUF_GET_PARAM(bufid, zb_zcl_device_callback_param_t);

	LOG_INF("%s id %hd", __func__, device_cb_param->device_cb_id);

	/* Set default response value. */
	device_cb_param->status = RET_OK;

	switch (device_cb_param->device_cb_id) {
	case ZB_ZCL_SET_ATTR_VALUE_CB_ID:
		cluster_id = device_cb_param->cb_param.
			     set_attr_value_param.cluster_id;
		attr_id = device_cb_param->cb_param.
			  set_attr_value_param.attr_id;

		if (cluster_id == ZB_ZCL_CLUSTER_ID_ON_OFF) {
			uint8_t value =
				device_cb_param->cb_param.set_attr_value_param
				.values.data8;

			LOG_INF("on/off attribute setting to %hd", value);
			if (attr_id == ZB_ZCL_ATTR_ON_OFF_ON_OFF_ID) {
				// on_off_set_value((zb_bool_t)value);
			}
		} else if (cluster_id == ZB_ZCL_CLUSTER_ID_LEVEL_CONTROL) {
		} else {
			/* Other clusters can be processed here */
			LOG_INF("Unhandled cluster attribute id: %d",
				cluster_id);
			device_cb_param->status = RET_NOT_IMPLEMENTED;
		}
		break;

	default:
		// if (zcl_scenes_cb(bufid) == ZB_FALSE) {
		// 	device_cb_param->status = RET_NOT_IMPLEMENTED;
		// }
		break;
	}

	LOG_INF("%s status: %hd", __func__, device_cb_param->status);
}

void zboss_signal_handler(zb_bufid_t bufid)
{
	zb_zdo_app_signal_hdr_t *signal_header = NULL;
	zb_zdo_app_signal_type_t signal = zb_get_app_signal(bufid, &signal_header);
	zb_ret_t err = RET_OK;

	/* Update network status LED but only for debug configuration */
	#ifdef CONFIG_LOG
	zigbee_led_status_update(bufid, ZIGBEE_NETWORK_STATE_LED);
	#endif /* CONFIG_LOG */

	/* Detect ZBOSS startup */
	switch (signal) {
	case ZB_ZDO_SIGNAL_SKIP_STARTUP:
		/* ZBOSS framework has started - schedule first weather check */
		err = ZB_SCHEDULE_APP_ALARM(check_weather,
					    0,
					    ZB_MILLISECONDS_TO_BEACON_INTERVAL(
						    WEATHER_CHECK_INITIAL_DELAY_MSEC));
		if (err) {
			LOG_ERR("Failed to schedule app alarm: %d", err);
		}
		break;
	default:
		break;
	}

	/* Let default signal handler process the signal*/
	ZB_ERROR_CHECK(zigbee_default_signal_handler(bufid));

	/*
	 * All callbacks should either reuse or free passed buffers.
	 * If bufid == 0, the buffer is invalid (not passed).
	 */
	if (bufid) {
		zb_buf_free(bufid);
	}
}

int main(void)
{
	#ifdef CONFIG_USB_DEVICE_STACK
	wait_for_console();
	#endif /* CONFIG_USB_DEVICE_STACK */

	LOG_INF("App init");

	register_factory_reset_button(FACTORY_RESET_BUTTON);
	gpio_init();
	weather_station_init();

	/* Register device context (endpoint) */
	ZB_AF_REGISTER_DEVICE_CTX(&weather_station_ctx);

	/* Register callback for handling ZCL commands. */
	ZB_ZCL_REGISTER_DEVICE_CB(zcl_device_cb);

	/* Init Basic and Identify attributes */
	mandatory_clusters_attr_init();

	/* Init measurements-related attributes */
	measurements_clusters_attr_init();

	/* Register callback to identify notifications */
	ZB_AF_SET_IDENTIFY_NOTIFICATION_HANDLER(WEATHER_STATION_ENDPOINT_NB, identify_callback);

	/* Enable Sleepy End Device behavior */
	zb_set_rx_on_when_idle(ZB_FALSE);
	if (IS_ENABLED(CONFIG_RAM_POWER_DOWN_LIBRARY)) {
		power_down_unused_ram();
	}

	/* Start Zigbee stack */
	zigbee_enable();

	dk_set_led_on(LED_POWER);

	return 0;
}
