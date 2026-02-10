#pragma once

#include <memory>
#include <vector>
#include <unordered_map>

#ifdef __cplusplus
extern "C"
{
#endif

#include <zboss_api.h>
#include <zboss_api_addons.h>

#ifdef __cplusplus
}
#endif

/* Delay for console initialization */
#define WAIT_FOR_CONSOLE_MSEC 100
#define WAIT_FOR_CONSOLE_DEADLINE_MSEC 500
/* Time of LED on state while blinking for identify mode */
#define IDENTIFY_LED_BLINK_TIME_MSEC 500

/* LED indicating that device successfully joined Zigbee network */
#define ZIGBEE_NETWORK_STATE_LED {{ if not (eq .Device.Board.NetworkStateLED "") -}}
{{ toButtonBit .Device.Board.NetworkStateLED }}
{{- else -}}LED_BLUE{{ end }}

/* LED used for device identification */
#define IDENTIFY_LED LED_RED

/* Button used to enter the Identify mode */
#define IDENTIFY_MODE_BUTTON {{ if not (eq .Device.Board.FactoryResetButton "") -}}
{{ toButtonBit .Device.Board.FactoryResetButton }}
{{- else -}}DK_BTN1_MSK{{ end }}

/* Button to start Factory Reset */
#define FACTORY_RESET_BUTTON IDENTIFY_MODE_BUTTON

#define MANUFACTURER_CODE ZB_ZCL_MANUF_CODE_INVALID
#define ZB_HA_DEVICE_VER 1

/* Manufacturer name (32 bytes). */
#define DEVICE_INIT_BASIC_MANUF_NAME "{{ .Device.General.Manufacturer }}"

/* Model number assigned by manufacturer (32-bytes long string). */
#define DEVICE_INIT_BASIC_MODEL_ID "{{ .Device.General.DeviceName }}"

/* First 8 bytes specify the date of manufacturer of the device
 * in ISO 8601 format (YYYYMMDD). The rest (8 bytes) are manufacturer specific.
 */
#define DEVICE_INIT_BASIC_DATE_CODE "{{ .GeneratedOn.Format `20060102` }}" // "20250923"

/* Describes the physical location of the device (16 bytes).
 * May be modified during commissioning process.
 */
#define DEVICE_INIT_BASIC_LOCATION_DESC ""
/* Describes the type of physical environment.
 * For possible values see section 3.2.2.2.10 of ZCL specification.
 */
#define DEVICE_INIT_BASIC_PH_ENV ZB_ZCL_BASIC_ENV_UNSPECIFIED

/* Zigbee Cluster Library 4.4.2.1.1: MeasuredValue = 100x temperature in degrees Celsius */
#define ZCL_TEMPERATURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER 100
/* Zigbee Cluster Library 4.5.2.2.1.1: MeasuredValue = 10x pressure in kPa */
#define ZCL_PRESSURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER 10
/* Zigbee Cluster Library 4.7.2.1.1: MeasuredValue = 100x water content in % */
#define ZCL_HUMIDITY_MEASUREMENT_MEASURED_VALUE_MULTIPLIER 100

typedef ZB_PACKED_PRE struct zb_af_simple_desc_s
{
    zb_uint8_t endpoint;                                                        /* Endpoint */
    zb_uint16_t app_profile_id;                                                 /* Application profile identifier */
    zb_uint16_t app_device_id;                                                  /* Application device identifier */
    zb_bitfield_t app_device_version : 4;                                       /* Application device version */
    zb_bitfield_t reserved : 4;                                                 /* Reserved */
    zb_uint8_t app_input_cluster_count;                                         /* Application input cluster count */
    zb_uint8_t app_output_cluster_count; /* Application output cluster count */ /* Application input and output cluster list */
    zb_uint16_t app_cluster_list;
} ZB_PACKED_STRUCT
    zb_af_simple_desc_t;

namespace zbhome
{
    namespace experimental
    {
        namespace zigbee
        {
            typedef struct
            {
                zb_int8_t zone_state;
                zb_int16_t zone_type;
                zb_int16_t zone_status;
                zb_uint8_t zone_id;
                zb_int64_t ias_cie_address;
                zb_int8_t cie_short_addr;
                zb_int16_t cie_ep;
                zb_uint8_t number_of_zone_sens_levels_supported;
                zb_uint8_t current_zone_sens_level;
            } zb_zcl_ias_zone_attrs_t;

            typedef struct
            {
                zb_uint16_t measure_value;
                zb_uint16_t min_measure_value;
                zb_uint16_t max_measure_value;
            } zb_zcl_water_content_attrs_t;

            typedef struct
            {
                zb_int16_t measure_value;
                zb_int16_t min_measure_value;
                zb_int16_t max_measure_value;
                zb_uint16_t tolerance;
            } zb_zcl_pressure_measurement_attrs_t;

            class Cluster
            {
            public:
                virtual const zb_uint16_t id() = 0;

                virtual zb_uint8_t role()
                {
                    return ZB_ZCL_CLUSTER_SERVER_ROLE;
                };

                virtual zb_uint8_t reporting_attr_count()
                {
                    return 0;
                };

                std::vector<zb_zcl_attr_t> attr_list;
                zb_zcl_cluster_desc_t cluster;
            };

            class BasicCluster : public Cluster
            {
            public:
                const zb_uint16_t id()
                {
                    return ZB_ZCL_CLUSTER_ID_BASIC;
                };

                zb_zcl_basic_attrs_ext_t attrs;

                // Specific attributes
                zb_bool_t device_enable_attr_list = 1U;
                zb_uint16_t cluster_revision_attr_list = ((zb_uint16_t)0x0003u);
            };

            class OnOffCluster : public Cluster
            {
            public:
                const zb_uint16_t id()
                {
                    return ZB_ZCL_CLUSTER_ID_ON_OFF;
                };

                zb_uint8_t role()
                {
                    return is_output ? ZB_ZCL_CLUSTER_SERVER_ROLE : ZB_ZCL_CLUSTER_CLIENT_ROLE;
                };

                zb_uint8_t reporting_attr_count()
                {
                    return is_output ? 1 : 0;
                };

                zb_zcl_on_off_attrs_t attrs;

                bool is_output;

                // Specific attributes
                zb_uint16_t cluster_revision_attr_list = ((zb_uint16_t)0x0002u);
            };

            class TemperatureCluster : public Cluster
            {
            public:
                const zb_uint16_t id()
                {
                    return ZB_ZCL_CLUSTER_ID_TEMP_MEASUREMENT;
                };

                zb_uint8_t reporting_attr_count()
                {
                    return 1;
                };

                zb_zcl_temp_measurement_attrs_t attrs;

                // Specific attributes
                zb_uint16_t cluster_revision_attr_list = ((zb_uint16_t)0x0003u);
            };

            // Humidity cluster can also be used for Leaf Wetness and Soil Moisture,
            // but for now - it is only humidity.
            //
            // This comes mainly to the defined ID and the ID used in cluster constructor.
            class HumidityCluster : public Cluster
            {
            public:
                const zb_uint16_t id()
                {
                    return ZB_ZCL_CLUSTER_ID_REL_HUMIDITY_MEASUREMENT;
                };
                zb_uint8_t reporting_attr_count()
                {
                    return 1;
                };

                zb_zcl_water_content_attrs_t attrs;

                zb_uint16_t cluster_revision_attr_list = ((zb_uint16_t)0x0002u);
            };

            class PressureCluster : public Cluster
            {
            public:
                const zb_uint16_t id()
                {
                    return ZB_ZCL_CLUSTER_ID_PRESSURE_MEASUREMENT;
                };

                zb_uint8_t reporting_attr_count()
                {
                    return 1;
                };

                zb_zcl_temp_measurement_attrs_t attrs;

                zb_uint16_t cluster_revision_attr_list = ((zb_uint16_t)0x0002u);
            };

            class IASZoneCluster : public Cluster
            {
            public:
                const zb_uint16_t id()
                {
                    return ZB_ZCL_CLUSTER_ID_IAS_ZONE;
                };

                zb_zcl_ias_zone_attrs_t attrs;

                // Specific values
                zb_uint8_t cie_addr_is_set_attr_list;
                zb_zcl_ias_zone_int_ctx_t int_ctx_attr_list;
                zb_uint16_t cluster_revision_attr_list = ((zb_uint16_t)0x0002u);
            };

            BasicCluster *basic_cluster();
            IASZoneCluster *ias_zone_cluster();
            OnOffCluster *on_off_cluster(bool is_output = true);
            TemperatureCluster *temperature_cluster(zb_int8_t min_measured, zb_int8_t max_measured, zb_int8_t tolerance);
            HumidityCluster *humidity_cluster(zb_int8_t min_measured, zb_int8_t max_measured, zb_int8_t tolerance);
            PressureCluster *pressure_cluster(zb_int8_t min_measured, zb_int8_t max_measured, zb_int8_t tolerance);

            // This classes exists because all this data needs to be held anyway,
            // and they are connected together.
            class Endpoint
            {
            public:
                zb_uint8_t endpoint_id()
                {
                    return ep.ep_id;
                };

                // zb_af_simple_desc_t desc_ep;
                zb_af_simple_desc_8_9_t desc_ep;
                std::vector<zb_zcl_reporting_info_t> reporting_info;
                std::vector<Cluster *> clusters;
                std::vector<zb_zcl_cluster_desc_t> zb_clusters;

                zb_af_endpoint_desc_t ep;
            };

            // ComplexEndpoint is endpoint has has at most 8 input & 9 output clusters.
            // class ComplexEndpoint : public Endpoint
            // {
            // public:
            // };

            Endpoint *endpoint(std::vector<Cluster *> clusters);
        }
    }
}