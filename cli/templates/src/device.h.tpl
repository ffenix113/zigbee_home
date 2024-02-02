#pragma once

#include "clusters.h"

/* Delay for console initialization */
#define WAIT_FOR_CONSOLE_MSEC 100
#define WAIT_FOR_CONSOLE_DEADLINE_MSEC 500

/* Weather check period */
#define WEATHER_CHECK_PERIOD_MSEC {{.Device.General.RunEvery.Milliseconds}}

/* Delay for first weather check */
#define CONFIG_FIRST_WEATHER_CHECK_DELAY_SECONDS 5
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

{{template "device_ctx" .Device.Sensors }}

/* Stores all cluster-related attributes */
static zb_device_ctx dev_ctx;

/* Attributes setup */
{{ range $i, $sensor := .Device.Sensors }}
	{{range $_, $cluster := $sensor.Clusters}}
		{{ render (clusterTpl $cluster.ID "attr_list") (clusterCtx $i $cluster)}}
	{{end}}
{{end}}

/* Clusters setup */
{{- $sensorsLen := len .Device.Sensors}}
{{- range $sensorIdx, $sensor := .Device.Sensors}}
ZB_HA_DECLARE_DEVICE_CLUSTER_LIST_EP_{{sum $sensorIdx 1}}(
	device_cluster_list_{{$sensorIdx}},
	{{- $clustersLen := len $sensor.Clusters}}
	{{- range $i, $cluster := $sensor.Clusters}}
	{{$cluster.CVarName}}_{{$sensorIdx}}_attr_list{{if not (isLast $i $clustersLen)}},{{end}}
	{{- end}}
	);
{{- end}}

/* Endpoint setup (single) */
{{- range $i, $sensor := .Device.Sensors}}
ZB_HA_DECLARE_DEVICE_EP(
	device_ep_{{$i}},
	{{sum $i 1}},
	{{$sensor.Clusters.ReportAttrCount}},
	device_cluster_list_{{$i}});
{{- end}}

/* Device context */
ZBOSS_DECLARE_DEVICE_CTX_{{len .Device.Sensors}}_EP(
	device_ctx,
	{{- range $i, $_ := .Device.Sensors}}
	device_ep_{{$i}}{{if not (isLast $i $sensorsLen)}},{{end}}
	{{- end}}
);

/* Manufacturer name (32 bytes). */
#define DEVICE_INIT_BASIC_MANUF_NAME      "FFenix113"

/* Model number assigned by manufacturer (32-bytes long string). */
#define DEVICE_INIT_BASIC_MODEL_ID        "dongle"

/* First 8 bytes specify the date of manufacturer of the device
 * in ISO 8601 format (YYYYMMDD). The rest (8 bytes) are manufacturer specific.
 */
#define DEVICE_INIT_BASIC_DATE_CODE       "{{ .GeneratedOn.Format "20060102" }}"

/* Describes the physical location of the device (16 bytes).
 * May be modified during commissioning process.
 */
#define DEVICE_INIT_BASIC_LOCATION_DESC   ""
/* Describes the type of physical environment.
 * For possible values see section 3.2.2.2.10 of ZCL specification.
 */
#define DEVICE_INIT_BASIC_PH_ENV          ZB_ZCL_BASIC_ENV_UNSPECIFIED

static void mandatory_clusters_attr_init(void)
{
	/* Basic cluster attributes */
	dev_ctx.basic_attr.zcl_version = ZB_ZCL_VERSION;
	dev_ctx.basic_attr.power_source = ZB_ZCL_BASIC_POWER_SOURCE_DC_SOURCE;
	// Extended attributes
	dev_ctx.basic_attr.app_version = 1;
	dev_ctx.basic_attr.stack_version = 1;
	dev_ctx.basic_attr.hw_version = 1;

	/* Use ZB_ZCL_SET_STRING_VAL to set strings, because the first byte
	 * should contain string length without trailing zero.
	 *
	 * For example "test" string will be encoded as:
	 *   [(0x4), 't', 'e', 's', 't']
	 */
	ZB_ZCL_SET_STRING_VAL(
		dev_ctx.basic_attr.mf_name,
		DEVICE_INIT_BASIC_MANUF_NAME,
		ZB_ZCL_STRING_CONST_SIZE(DEVICE_INIT_BASIC_MANUF_NAME));

	ZB_ZCL_SET_STRING_VAL(
		dev_ctx.basic_attr.model_id,
		DEVICE_INIT_BASIC_MODEL_ID,
		ZB_ZCL_STRING_CONST_SIZE(DEVICE_INIT_BASIC_MODEL_ID));

	ZB_ZCL_SET_STRING_VAL(
		dev_ctx.basic_attr.date_code,
		DEVICE_INIT_BASIC_DATE_CODE,
		ZB_ZCL_STRING_CONST_SIZE(DEVICE_INIT_BASIC_DATE_CODE));

	ZB_ZCL_SET_STRING_VAL(
		dev_ctx.basic_attr.location_id,
		DEVICE_INIT_BASIC_LOCATION_DESC,
		ZB_ZCL_STRING_CONST_SIZE(DEVICE_INIT_BASIC_LOCATION_DESC));

	dev_ctx.basic_attr.ph_env = DEVICE_INIT_BASIC_PH_ENV;

	/* Identify cluster attributes data. */
	// dev_ctx.identify_attr.identify_time =
	// 	ZB_ZCL_IDENTIFY_IDENTIFY_TIME_DEFAULT_VALUE;

	/* Identify cluster attributes */
	// dev_ctx.identify_attr.identify_time = ZB_ZCL_IDENTIFY_IDENTIFY_TIME_DEFAULT_VALUE;
}

static void measurements_clusters_attr_init(void)
{
	{{- range $i, $sensor := .Device.Sensors}}
	{{- range $sensor.Clusters}}
	{{- maybeRender (clusterTpl .ID "attr_init") (clusterCtx $i .)}}
	{{- end}}
	{{- end}}
}
