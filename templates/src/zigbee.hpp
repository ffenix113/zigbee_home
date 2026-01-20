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
    namespace zigbee
    {
        namespace experimental
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