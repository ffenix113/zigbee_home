{{ define "define_cluster_list" }}
#include <zboss_api.h>

/* Number chosen for the single endpoint provided by device */
#define DEVICE_ENDPOINT_NB 1

/* Temperature sensor device version */
#define ZB_HA_DEVICE_VER_TEMPERATURE_SENSOR     0
/* Basic, identify, temperature, pressure, humidity, on/off */
// Have basic + identify as always present, and other as optional.
#define ZB_HA_DEVICE_IN_CLUSTER_NUM    {{ sum 2 .Clusters.ReportAttrCount }}
/* Identify */ // Not for generated device for now.
#define ZB_HA_DEVICE_OUT_CLUSTER_NUM   0

/* Temperature, pressure, humidity, on/off */
#define ZB_HA_DEVICE_REPORT_ATTR_COUNT {{ .Clusters.ReportAttrCount }}

#define MANUFACTURER_CODE ZB_ZCL_MANUF_CODE_INVALID

#define ZB_HA_DECLARE_DEVICE_CLUSTER_LIST(						\
		cluster_list_name,								\
		basic_attr_list,								\
		{{- $clustersLen := len .Clusters}}
        {{- range $i, $cluster := .Clusters }}
        {{ $cluster.CVarName }}_attr_list{{if not (isLast $i $clustersLen)}},{{end}} \
        {{- end }}
		)										\
	zb_zcl_cluster_desc_t cluster_list_name[] =						\
	{											\
		{{- range .Clusters}}
		ZB_ZCL_CLUSTER_DESC(								\
			{{.ID}},						\
			ZB_ZCL_ARRAY_SIZE({{.CVarName}}_attr_list, zb_zcl_attr_t),			\
			({{.CVarName}}_attr_list),							\
			ZB_ZCL_CLUSTER_SERVER_ROLE,	 /* For now let's say all are server role. Later this can be cluster paramter */					\
			MANUFACTURER_CODE						\
			),									\
		{{- end}}
	}

#define ZB_ZCL_DECLARE_DEVICE_DESC(						\
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
			{{- range .Clusters}}
			{{.ID}}, \
			{{- end}}
		}									\
	}

#define ZB_HA_DECLARE_DEVICE_EP(ep_name, ep_id, cluster_list)				\
	ZB_ZCL_DECLARE_DEVICE_DESC(							\
		ep_name,									\
		ep_id,										\
		ZB_HA_DEVICE_IN_CLUSTER_NUM,						\
		ZB_HA_DEVICE_OUT_CLUSTER_NUM);						\
	ZBOSS_DEVICE_DECLARE_REPORTING_CTX(							\
		reporting_info##ep_name,							\
		ZB_HA_DEVICE_REPORT_ATTR_COUNT);					\
	ZB_AF_DECLARE_ENDPOINT_DESC(								\
		ep_name,									\
		ep_id,										\
		ZB_AF_HA_PROFILE_ID,								\
		0,										\
		NULL,										\
		ZB_ZCL_ARRAY_SIZE(cluster_list, zb_zcl_cluster_desc_t),				\
		cluster_list,									\
		(zb_af_simple_desc_1_1_t *)&simple_desc_##ep_name,				\
		ZB_HA_DEVICE_REPORT_ATTR_COUNT, reporting_info##ep_name, 0, NULL)
{{ end }}
{{- template "define_cluster_list" index .Device.Endpoints 0 }}