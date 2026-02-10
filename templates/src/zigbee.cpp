#include <memory>
#include <vector>

#ifdef __cplusplus
extern "C"
{
#endif

#include <zboss_api.h>
#include <zboss_api_addons.h>

#ifdef __cplusplus
}
#endif

#include "zephyr/logging/log.h"

#include "zigbee.hpp"

LOG_MODULE_REGISTER(zigbee, LOG_LEVEL_DBG);

namespace zbhome
{
    namespace experimental
    {
        // static zb_af_device_ctx_t device_ctx;

        namespace zigbee
        {

            zb_uint8_t next_endpoint_id()
            {
                static zb_uint8_t ep_id = 0;
                ep_id++;

                return ep_id;
            };

            BasicCluster *basic_cluster()
            {
                // auto basic = std::make_unique<BasicCluster>();
                auto cluster = new BasicCluster();

                cluster->attrs.zcl_version = ZB_ZCL_VERSION;
                cluster->attrs.power_source = ZB_ZCL_BASIC_POWER_SOURCE_DC_SOURCE;
                cluster->attrs.ph_env = DEVICE_INIT_BASIC_PH_ENV;
                // Extended attributes
                cluster->attrs.app_version = 1;
                cluster->attrs.stack_version = 1;
                cluster->attrs.hw_version = 1;

                /* Use ZB_ZCL_SET_STRING_VAL to set strings, because the first byte
                 * should contain string length without trailing zero.
                 *
                 * For example "test" string will be encoded as:
                 *   [(0x4), 't', 'e', 's', 't']
                 */
                ZB_ZCL_SET_STRING_VAL(
                    cluster->attrs.mf_name,
                    DEVICE_INIT_BASIC_MANUF_NAME,
                    ZB_ZCL_STRING_CONST_SIZE(DEVICE_INIT_BASIC_MANUF_NAME));

                ZB_ZCL_SET_STRING_VAL(
                    cluster->attrs.model_id,
                    DEVICE_INIT_BASIC_MODEL_ID,
                    ZB_ZCL_STRING_CONST_SIZE(DEVICE_INIT_BASIC_MODEL_ID));

                ZB_ZCL_SET_STRING_VAL(
                    cluster->attrs.date_code,
                    DEVICE_INIT_BASIC_DATE_CODE,
                    ZB_ZCL_STRING_CONST_SIZE(DEVICE_INIT_BASIC_DATE_CODE));

                ZB_ZCL_SET_STRING_VAL(
                    cluster->attrs.location_id,
                    DEVICE_INIT_BASIC_LOCATION_DESC,
                    ZB_ZCL_STRING_CONST_SIZE(DEVICE_INIT_BASIC_LOCATION_DESC));

                // ZB_ZCL_DECLARE_BASIC_ATTRIB_LIST_EXT(
                //     attr_list,
                //     &cluster->attrs.zcl_version,
                //     &cluster->attrs.app_version,
                //     &cluster->attrs.stack_version,
                //     &cluster->attrs.hw_version,
                //     cluster->attrs.mf_name,
                //     cluster->attrs.model_id,
                //     cluster->attrs.date_code,
                //     &cluster->attrs.power_source,
                //     cluster->attrs.location_id,
                //     &cluster->attrs.ph_env,
                //     cluster->attrs.sw_ver);

                cluster->attr_list = {
                    {0xfffdU, 0x21U, 0x01U, 0xFFFFU, (void *)&(cluster->cluster_revision_attr_list)},
                    {ZB_ZCL_ATTR_BASIC_ZCL_VERSION_ID, 0x20U, 0x01U, (0xFFFFU), (void *)(&cluster->attrs.zcl_version)},
                    {ZB_ZCL_ATTR_BASIC_APPLICATION_VERSION_ID, 0x20U, 0x01U, (0xFFFFU), (void *)(&cluster->attrs.app_version)},
                    {ZB_ZCL_ATTR_BASIC_STACK_VERSION_ID, 0x20U, 0x01U, (0xFFFFU), (void *)(&cluster->attrs.stack_version)},
                    {ZB_ZCL_ATTR_BASIC_HW_VERSION_ID, 0x20U, 0x01U, (0xFFFFU), (void *)(&cluster->attrs.hw_version)},
                    {ZB_ZCL_ATTR_BASIC_MANUFACTURER_NAME_ID, 0x42U, 0x01U, (0xFFFFU), (void *)(cluster->attrs.mf_name)},
                    {ZB_ZCL_ATTR_BASIC_MODEL_IDENTIFIER_ID, 0x42U, 0x01U, (0xFFFFU), (void *)(cluster->attrs.model_id)},
                    {ZB_ZCL_ATTR_BASIC_DATE_CODE_ID, 0x42U, 0x01U, (0xFFFFU), (void *)(cluster->attrs.date_code)},
                    {ZB_ZCL_ATTR_BASIC_POWER_SOURCE_ID, 0x30U, 0x01U, (0xFFFFU), (void *)(&cluster->attrs.power_source)},
                    {ZB_ZCL_ATTR_BASIC_SW_BUILD_ID, 0x42U, 0x01U, (0xFFFFU), (void *)(cluster->attrs.sw_ver)},
                    {ZB_ZCL_ATTR_BASIC_DEVICE_ENABLED_ID, 0x10U, (0x01U | 0x02U), (0xFFFFU), (void *)&(cluster->device_enable_attr_list)},
                    {ZB_ZCL_ATTR_BASIC_LOCATION_DESCRIPTION_ID, 0x42U, (0x01U | 0x02U), (0xFFFFU), (void *)(cluster->attrs.location_id)},
                    {ZB_ZCL_ATTR_BASIC_PHYSICAL_ENVIRONMENT_ID, 0x30U, (0x01U | 0x02U), (0xFFFFU), (void *)(&cluster->attrs.ph_env)},
                    {(zb_uint16_t)(-1), 0, 0, 0xFFFFU, __null}};

                cluster->cluster = ZB_ZCL_CLUSTER_DESC(
                    ZB_ZCL_CLUSTER_ID_BASIC,
                    zb_uint8_t(cluster->attr_list.size()), // ZB_ZCL_ARRAY_SIZE(attr_list, zb_zcl_attr_t),
                    (&cluster->attr_list[0]),
                    cluster->role(),
                    MANUFACTURER_CODE);

                return cluster;
            };

            OnOffCluster *on_off_cluster(bool is_output)
            {
                auto cluster = new OnOffCluster();
                cluster->is_output = is_output;

                cluster->attrs.on_off = ZB_ZCL_ON_OFF_IS_OFF;

                // ZB_ZCL_DECLARE_ON_OFF_ATTRIB_LIST(
                //     attr_list,
                //     &cluster->attrs.on_off);

                cluster->attr_list = {
                    {0xfffdU, 0x21U, 0x01U, 0xFFFFU, (void *)&(cluster->cluster_revision_attr_list)},
                    {ZB_ZCL_ATTR_ON_OFF_ON_OFF_ID, 0x10U, 0x01U | 0x04U | 0x10U, (0xFFFFU), (void *)(&cluster->attrs.on_off)},
                    {(zb_uint16_t)(-1), 0, 0, 0xFFFFU, __null}};

                if (is_output)
                {
                    cluster->cluster = ZB_ZCL_CLUSTER_DESC(
                        ZB_ZCL_CLUSTER_ID_ON_OFF,
                        zb_uint8_t(cluster->attr_list.size()),
                        (&cluster->attr_list[0]),
                        cluster->role(),
                        MANUFACTURER_CODE);
                }
                else
                {
                    cluster->cluster = ZB_ZCL_CLUSTER_DESC(
                        ZB_ZCL_CLUSTER_ID_ON_OFF,
                        zb_uint8_t(cluster->attr_list.size()),
                        (&cluster->attr_list[0]),
                        cluster->role(),
                        MANUFACTURER_CODE);
                }

                return cluster;
            };

            IASZoneCluster *ias_zone_cluster()
            {
                auto cluster = new IASZoneCluster();

                cluster->attrs.zone_state = ZB_ZCL_IAS_ZONE_ZONESTATE_DEF_VALUE;
                cluster->attrs.zone_type = ZB_ZCL_IAS_ZONE_ZONETYPE_CONTACT_SWITCH;
                // Set status to include Restore, to specify that
                // contact sensor does know when it is closed & open.
                cluster->attrs.zone_status = ZB_ZCL_IAS_ZONE_ZONE_STATUS_RESTORE;
                cluster->attrs.number_of_zone_sens_levels_supported = ZB_ZCL_IAS_ZONE_NUMBER_OF_ZONE_SENSITIVITY_LEVELS_SUPPORTED_DEFAULT_VALUE;
                cluster->attrs.current_zone_sens_level = ZB_ZCL_IAS_ZONE_CURRENT_ZONE_SENSITIVITY_LEVEL_DEFAULT_VALUE;
                cluster->attrs.zone_id = ZB_ZCL_IAS_ZONEID_ID_DEF_VALUE;

                cluster->attr_list = {
                    {0xfffdU, 0x21U, 0x01U, 0xFFFFU, (void *)&(cluster->cluster_revision_attr_list)},
                    {ZB_ZCL_ATTR_IAS_ZONE_ZONESTATE_ID, 0x30U, 0x01U, (0xFFFFU), (void *)(&cluster->attrs.zone_state)},
                    {ZB_ZCL_ATTR_IAS_ZONE_ZONETYPE_ID, 0x31U, 0x01U, (0xFFFFU), (void *)(&cluster->attrs.zone_type)},
                    {ZB_ZCL_ATTR_IAS_ZONE_ZONESTATUS_ID, 0x19U, 0x01U | 0x04U, (0xFFFFU), (void *)(&cluster->attrs.zone_status)},
                    {ZB_ZCL_ATTR_IAS_ZONE_IAS_CIE_ADDRESS_ID, 0xf0U, (0x01U | 0x02U), (0xFFFFU), (void *)(&cluster->attrs.ias_cie_address)},
                    {ZB_ZCL_ATTR_IAS_ZONE_ZONEID_ID, 0x20U, 0x01U, (0xFFFFU), (void *)(&cluster->attrs.zone_id)},
                    {ZB_ZCL_ATTR_IAS_ZONE_NUMBER_OF_ZONE_SENSITIVITY_LEVELS_SUPPORTED_ID, 0x20U, 0x01U, (0xFFFFU), (void *)(&cluster->attrs.number_of_zone_sens_levels_supported)},
                    {ZB_ZCL_ATTR_IAS_ZONE_CURRENT_ZONE_SENSITIVITY_LEVEL_ID, 0x20U, (0x01U | 0x02U), (0xFFFFU), (void *)(&cluster->attrs.current_zone_sens_level)},
                    {ZB_ZCL_ATTR_IAS_ZONE_INT_CTX_ID, 0x00U, 0x40U, (0xFFFFU), (void *)&(cluster->int_ctx_attr_list)},
                    {ZB_ZCL_ATTR_CUSTOM_CIE_SHORT_ADDR, 0x21U, 0x40U, (0xFFFFU), (void *)(&cluster->attrs.cie_short_addr)},
                    {ZB_ZCL_ATTR_CUSTOM_CIE_EP, 0x20U, 0x40U, (0xFFFFU), (void *)(&cluster->attrs.cie_ep)},
                    {ZB_ZCL_ATTR_CUSTOM_CIE_ADDR_IS_SET, 0x20U, 0x40U, (0xFFFFU), (void *)&(cluster->cie_addr_is_set_attr_list)},
                    {(zb_uint16_t)(-1), 0, 0, 0xFFFFU, __null}};

                cluster->cluster = ZB_ZCL_CLUSTER_DESC(
                    ZB_ZCL_CLUSTER_ID_IAS_ZONE,
                    zb_uint8_t(cluster->attr_list.size()),
                    (&cluster->attr_list[0]),
                    cluster->role(),
                    MANUFACTURER_CODE);

                return cluster;
            };

            // Temperature
            TemperatureCluster *temperature_cluster(zb_int8_t min_measured, zb_int8_t max_measured, zb_int8_t tolerance)
            {
                auto cluster = new TemperatureCluster();

                cluster->attrs.measure_value = ZB_ZCL_ATTR_TEMP_MEASUREMENT_MIN_VALUE_MIN_VALUE;
                cluster->attrs.min_measure_value = min_measured * ZCL_TEMPERATURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER;
                cluster->attrs.max_measure_value = max_measured * ZCL_TEMPERATURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER;
                cluster->attrs.tolerance = tolerance * ZCL_TEMPERATURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER;

                // ZB_ZCL_DECLARE_TEMP_MEASUREMENT_ATTRIB_LIST(
                //     attr_list,
                //     &cluster->attrs.measure_value,
                //     &cluster->attrs.min_measure_value,
                //     &cluster->attrs.max_measure_value,
                //     &cluster->attrs.tolerance);

                cluster->attr_list = {
                    {0xfffdU, 0x21U, 0x01U, 0xFFFFU, (void *)&(cluster->cluster_revision_attr_list)},
                    {ZB_ZCL_ATTR_TEMP_MEASUREMENT_VALUE_ID, 0x29U, 0x01U | 0x04U, (0xFFFFU), (void *)(&cluster->attrs.measure_value)},
                    {ZB_ZCL_ATTR_TEMP_MEASUREMENT_MIN_VALUE_ID, 0x29U, 0x01U, (0xFFFFU), (void *)(&cluster->attrs.min_measure_value)},
                    {ZB_ZCL_ATTR_TEMP_MEASUREMENT_MAX_VALUE_ID, 0x29U, 0x01U, (0xFFFFU), (void *)(&cluster->attrs.max_measure_value)},
                    {ZB_ZCL_ATTR_TEMP_MEASUREMENT_TOLERANCE_ID, 0x21U, 0x01U, (0xFFFFU), (void *)(&cluster->attrs.tolerance)},
                    {(zb_uint16_t)(-1), 0, 0, 0xFFFFU, __null}};

                cluster->cluster = ZB_ZCL_CLUSTER_DESC(
                    ZB_ZCL_CLUSTER_ID_TEMP_MEASUREMENT,
                    zb_uint8_t(cluster->attr_list.size()),
                    (&cluster->attr_list[0]),
                    cluster->role(),
                    MANUFACTURER_CODE);

                return cluster;
            };

            HumidityCluster *humidity_cluster(zb_int8_t min_measured, zb_int8_t max_measured, zb_int8_t tolerance)
            {
                auto cluster = new HumidityCluster();

                cluster->attrs.measure_value = ZB_ZCL_ATTR_REL_HUMIDITY_MEASUREMENT_MIN_VALUE_MIN_VALUE;
                cluster->attrs.min_measure_value = min_measured * ZCL_HUMIDITY_MEASUREMENT_MEASURED_VALUE_MULTIPLIER;
                cluster->attrs.max_measure_value = max_measured * ZCL_HUMIDITY_MEASUREMENT_MEASURED_VALUE_MULTIPLIER;

                // ZB_ZCL_DECLARE_REL_HUMIDITY_MEASUREMENT_ATTRIB_LIST(
                //     attr_list,
                //     &cluster->attrs.measure_value,
                //     &cluster->attrs.min_measure_value,
                //     &cluster->attrs.max_measure_value);

                cluster->attr_list = {
                    {0xfffdU, 0x21U, 0x01U, 0xFFFFU, (void *)&(cluster->cluster_revision_attr_list)},
                    {ZB_ZCL_ATTR_REL_HUMIDITY_MEASUREMENT_VALUE_ID, 0x21U, 0x01U | 0x04U, (0xFFFFU), (void *)(&cluster->attrs.measure_value)},
                    {ZB_ZCL_ATTR_REL_HUMIDITY_MEASUREMENT_MIN_VALUE_ID, 0x21U, 0x01U, (0xFFFFU), (void *)(&cluster->attrs.min_measure_value)},
                    {ZB_ZCL_ATTR_REL_HUMIDITY_MEASUREMENT_MAX_VALUE_ID, 0x21U, 0x01U, (0xFFFFU), (void *)(&cluster->attrs.max_measure_value)},
                    {(zb_uint16_t)(-1), 0, 0, 0xFFFFU, __null}};

                cluster->cluster = ZB_ZCL_CLUSTER_DESC(
                    ZB_ZCL_CLUSTER_ID_REL_HUMIDITY_MEASUREMENT,
                    zb_uint8_t(cluster->attr_list.size()),
                    (&cluster->attr_list[0]),
                    cluster->role(),
                    MANUFACTURER_CODE);

                return cluster;
            };

            PressureCluster *pressure_cluster(zb_int8_t min_measured, zb_int8_t max_measured, zb_int8_t tolerance)
            {
                auto cluster = new PressureCluster();

                cluster->attrs.measure_value = ZB_ZCL_ATTR_PRESSURE_MEASUREMENT_MIN_VALUE_MIN_VALUE;
                cluster->attrs.min_measure_value = min_measured * ZCL_PRESSURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER;
                cluster->attrs.max_measure_value = max_measured * ZCL_PRESSURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER;

                // ZB_ZCL_DECLARE_PRESSURE_MEASUREMENT_ATTRIB_LIST(
                //     attr_list,
                //     &cluster->attrs.measure_value,
                //     &cluster->attrs.min_measure_value,
                //     &cluster->attrs.max_measure_value,
                //     &cluster->attrs.tolerance);

                cluster->attr_list = {
                    {0xfffdU, 0x21U, 0x01U, 0xFFFFU, (void *)&(cluster->cluster_revision_attr_list)},
                    {ZB_ZCL_ATTR_PRESSURE_MEASUREMENT_VALUE_ID, 0x29U, 0x01U | 0x04U, (0xFFFFU), (void *)(&cluster->attrs.measure_value)},
                    {ZB_ZCL_ATTR_PRESSURE_MEASUREMENT_MIN_VALUE_ID, 0x29U, 0x01U, (0xFFFFU), (void *)(&cluster->attrs.min_measure_value)},
                    {ZB_ZCL_ATTR_PRESSURE_MEASUREMENT_MAX_VALUE_ID, 0x29U, 0x01U, (0xFFFFU), (void *)(&cluster->attrs.max_measure_value)},
                    {ZB_ZCL_ATTR_PRESSURE_MEASUREMENT_TOLERANCE_ID, 0x21U, 0x01U, (0xFFFFU), (void *)(&cluster->attrs.tolerance)},
                    {(zb_uint16_t)(-1), 0, 0, 0xFFFFU, __null}};

                cluster->cluster = ZB_ZCL_CLUSTER_DESC(
                    ZB_ZCL_CLUSTER_ID_PRESSURE_MEASUREMENT,
                    zb_uint8_t(cluster->attr_list.size()),
                    (&cluster->attr_list[0]),
                    cluster->role(),
                    MANUFACTURER_CODE);

                return cluster;
            };

            // Endpoint constructor

            Endpoint *endpoint(std::vector<Cluster *> clusters)
            {
                auto ep = new Endpoint();
                auto endpoint_id = next_endpoint_id();

                zb_uint8_t num_input_clusters = 0;
                zb_uint8_t num_output_clusters = 0;
                zb_uint8_t num_reporting = 0;

                ep->clusters = std::move(clusters);

                ep->zb_clusters.reserve(ep->clusters.size());

                for (auto &cluster : ep->clusters)
                {
                    ep->zb_clusters.push_back(cluster->cluster);

                    if (cluster->role() & ZB_ZCL_CLUSTER_SERVER_ROLE)
                    {
                        num_input_clusters++;
                    };

                    if (cluster->role() & ZB_ZCL_CLUSTER_CLIENT_ROLE)
                    {
                        num_output_clusters++;
                    };

                    num_reporting += cluster->reporting_attr_count();
                }

                ep->desc_ep =
                    {
                        .endpoint = endpoint_id, // Static endpoint id for basic cluster endpoint.
                        .app_profile_id = ZB_AF_HA_PROFILE_ID,
                        .app_device_id = ZB_HA_TEMPERATURE_SENSOR_DEVICE_ID, /*This values are present as initial ones.*/
                        .app_device_version = ZB_HA_DEVICE_VER,              /*This values are present as initial ones.*/
                        .reserved = 0,                                       // always 0
                        .app_input_cluster_count = num_input_clusters,       // Input clusters
                        .app_output_cluster_count = num_output_clusters,     // Output clusters
                        // .app_cluster_list = {ep->cluster_ids[0]},
                    };

                zb_uint8_t cluster_list_idx = 0;
                for (auto &cluster_id : ep->clusters)
                {
                    ep->desc_ep.app_cluster_list[cluster_list_idx++] = cluster_id->id();
                }

                // ZBOSS_DEVICE_DECLARE_REPORTING_CTX(
                //     reporting_info, // FIXME: When this function returns - this will be lost. Not ideal.
                //     num_reporting);

                ep->reporting_info.resize(num_reporting);

                ZB_AF_DECLARE_ENDPOINT_DESC(
                    zb_ep,
                    endpoint_id,
                    ZB_AF_HA_PROFILE_ID,
                    0,
                    NULL,
                    zb_uint8_t(ep->zb_clusters.size()),
                    ep->zb_clusters.data(),
                    (zb_af_simple_desc_1_1_t *)&ep->desc_ep,
                    zb_uint8_t(ep->reporting_info.size()),
                    &ep->reporting_info[0],
                    0, NULL);
                ep->ep = zb_ep;

                LOG_DBG("created endpoint %d with %d cluster(s) and %d reporting", endpoint_id, ep->clusters.size(), ep->reporting_info.size());

                return ep;
            };
        }
    }
}