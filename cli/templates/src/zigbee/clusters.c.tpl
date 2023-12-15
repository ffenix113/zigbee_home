{{ define "defin_cluster_list" }}
#define ZB_HA_DECLARE_DEVICE_CLUSTER_LIST(						\
		cluster_list_name,								\
		basic_attr_list,								\
        {{ range .Clusters }}
        {{ clusterToCVar .}}_attr_list, \
        {{ end }}
		identify_client_attr_list,							\
		identify_server_attr_list,							\
		temperature_measurement_attr_list,						\
		pressure_measurement_attr_list,							\
		humidity_measurement_attr_list,							\
		on_off_attr_list							\
		)										\
	zb_zcl_cluster_desc_t cluster_list_name[] =						\
	{											\
		ZB_ZCL_CLUSTER_DESC(								\
			ZB_ZCL_CLUSTER_ID_BASIC,						\
			ZB_ZCL_ARRAY_SIZE(basic_attr_list, zb_zcl_attr_t),			\
			(basic_attr_list),							\
			ZB_ZCL_CLUSTER_SERVER_ROLE,						\
			MANUFACTURER_CODE						\
			),									\
		ZB_ZCL_CLUSTER_DESC(								\
			ZB_ZCL_CLUSTER_ID_IDENTIFY,						\
			ZB_ZCL_ARRAY_SIZE(identify_server_attr_list, zb_zcl_attr_t),		\
			(identify_server_attr_list),						\
			ZB_ZCL_CLUSTER_SERVER_ROLE,						\
			MANUFACTURER_CODE						\
			),									\
		ZB_ZCL_CLUSTER_DESC(								\
			ZB_ZCL_CLUSTER_ID_TEMP_MEASUREMENT,					\
			ZB_ZCL_ARRAY_SIZE(temperature_measurement_attr_list, zb_zcl_attr_t),	\
			(temperature_measurement_attr_list),					\
			ZB_ZCL_CLUSTER_SERVER_ROLE,						\
			MANUFACTURER_CODE						\
			),									\
		ZB_ZCL_CLUSTER_DESC(								\
			ZB_ZCL_CLUSTER_ID_PRESSURE_MEASUREMENT,					\
			ZB_ZCL_ARRAY_SIZE(pressure_measurement_attr_list, zb_zcl_attr_t),	\
			(pressure_measurement_attr_list),					\
			ZB_ZCL_CLUSTER_SERVER_ROLE,						\
			MANUFACTURER_CODE						\
			),									\
		ZB_ZCL_CLUSTER_DESC(								\
			ZB_ZCL_CLUSTER_ID_REL_HUMIDITY_MEASUREMENT,				\
			ZB_ZCL_ARRAY_SIZE(humidity_measurement_attr_list, zb_zcl_attr_t),	\
			(humidity_measurement_attr_list),					\
			ZB_ZCL_CLUSTER_SERVER_ROLE,						\
			MANUFACTURER_CODE						\
			),									\
		ZB_ZCL_CLUSTER_DESC(								\
			ZB_ZCL_CLUSTER_ID_ON_OFF,				\
			ZB_ZCL_ARRAY_SIZE(on_off_attr_list, zb_zcl_attr_t),	\
			(on_off_attr_list),					\
			ZB_ZCL_CLUSTER_SERVER_ROLE,						\
			MANUFACTURER_CODE						\
			),									\
		ZB_ZCL_CLUSTER_DESC(								\
			ZB_ZCL_CLUSTER_ID_IDENTIFY,						\
			ZB_ZCL_ARRAY_SIZE(identify_client_attr_list, zb_zcl_attr_t),		\
			(identify_client_attr_list),						\
			ZB_ZCL_CLUSTER_CLIENT_ROLE,						\
			MANUFACTURER_CODE						\
			),									\
	}

#define ZB_ZCL_DECLARE_WEATHER_STATION_DESC(						\
		ep_name,								\
		ep_id,									\
		in_clust_num,								\
		out_clust_num)								\
	ZB_DECLARE_SIMPLE_DESC(in_clust_num, out_clust_num);				\
	ZB_AF_SIMPLE_DESC_TYPE(in_clust_num, out_clust_num) simple_desc_##ep_name =	\
	{										\
		ep_id,									\
		ZB_AF_HA_PROFILE_ID,							\
		ZB_HA_TEMPERATURE_SENSOR_DEVICE_ID,					\
		ZB_HA_DEVICE_VER_TEMPERATURE_SENSOR,					\
		0,									\
		in_clust_num,								\
		out_clust_num,								\
		{									\
			ZB_ZCL_CLUSTER_ID_BASIC,					\
			ZB_ZCL_CLUSTER_ID_IDENTIFY,					\
			ZB_ZCL_CLUSTER_ID_TEMP_MEASUREMENT,				\
			ZB_ZCL_CLUSTER_ID_PRESSURE_MEASUREMENT,				\
			ZB_ZCL_CLUSTER_ID_REL_HUMIDITY_MEASUREMENT,			\
			ZB_ZCL_CLUSTER_ID_ON_OFF,							\
			ZB_ZCL_CLUSTER_ID_IDENTIFY,					\
		}									\
	}

#define ZB_HA_DECLARE_WEATHER_STATION_EP(ep_name, ep_id, cluster_list)				\
	ZB_ZCL_DECLARE_WEATHER_STATION_DESC(							\
		ep_name,									\
		ep_id,										\
		ZB_HA_WEATHER_STATION_IN_CLUSTER_NUM,						\
		ZB_HA_WEATHER_STATION_OUT_CLUSTER_NUM);						\
	ZBOSS_DEVICE_DECLARE_REPORTING_CTX(							\
		reporting_info##ep_name,							\
		ZB_HA_WEATHER_STATION_REPORT_ATTR_COUNT);					\
	ZB_AF_DECLARE_ENDPOINT_DESC(								\
		ep_name,									\
		ep_id,										\
		ZB_AF_HA_PROFILE_ID,								\
		0,										\
		NULL,										\
		ZB_ZCL_ARRAY_SIZE(cluster_list, zb_zcl_cluster_desc_t),				\
		cluster_list,									\
		(zb_af_simple_desc_1_1_t *)&simple_desc_##ep_name,				\
		ZB_HA_WEATHER_STATION_REPORT_ATTR_COUNT, reporting_info##ep_name, 0, NULL)
{{ end }}