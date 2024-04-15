#pragma once

#include <zboss_api.h>

#define MANUFACTURER_CODE ZB_ZCL_MANUF_CODE_INVALID

/* Temperature sensor device version */
#define ZB_HA_DEVICE_VER_TEMPERATURE_SENSOR     0

/* Zigbee Cluster Library 4.4.2.2.1.1: MeasuredValue = 100x temperature in degrees Celsius */
#define ZCL_TEMPERATURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER 100
/* Zigbee Cluster Library 4.5.2.2.1.1: MeasuredValue = 10x pressure in kPa */
#define ZCL_PRESSURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER 10
/* Zigbee Cluster Library 4.7.2.1.1: MeasuredValue = 100x water content in % */
#define ZCL_HUMIDITY_MEASUREMENT_MEASURED_VALUE_MULTIPLIER 100

{{ range $i, $cluster := .Device.Sensors.UniqueClusters }}
{{ maybeRender (clusterTpl $cluster.ID "defines") (clusterCtx $i $cluster)}}
{{end}}

{{- range $i, $sensor := .Device.Sensors }}
{{ $endpointID := (sum $i 1)}}
{{ $inClustersNum := $sensor.Clusters.Servers }}
{{if eq $i 0 }}{{ $inClustersNum = sum $inClustersNum 1 }}{{end}}
// Define a cluster array for a single endpoint
#define ZB_HA_DECLARE_DEVICE_CLUSTER_LIST_EP_{{$endpointID}}(						\
		cluster_list_name,								\
		{{- $clustersLen := len $sensor.Clusters}}
        {{- range $i, $cluster := $sensor.Clusters }}
        {{ $cluster.CVarName }}_attr_list{{if not (isLast $i $clustersLen)}},{{end}} \
        {{- end }}
		)										\
	zb_zcl_cluster_desc_t cluster_list_name[] =						\
	{											\
		{{- range $sensor.Clusters}}
		ZB_ZCL_CLUSTER_DESC(								\
			{{.ID.ToZCL}},						\
			ZB_ZCL_ARRAY_SIZE({{.CVarName}}_attr_list, zb_zcl_attr_t),			\
			({{.CVarName}}_attr_list),							\
			{{.Side.String}},	\
			MANUFACTURER_CODE						\
			),									\
		{{- end}}
	}
{{- end }}

{{- range $i, $sensor := .Device.Sensors}}
{{ $endpointID := (sum $i 1)}}
{{ $inClusterNum := $sensor.Clusters.Servers}}
{{ $outClusterNum := $sensor.Clusters.Clients}}
// Define an endpoint information(num of input&output, cluster IDs)
#define ZB_ZCL_DECLARE_DEVICE_DESC_EP_{{$endpointID}}(ep_name)								\
	ZB_DECLARE_SIMPLE_DESC_VA({{$inClusterNum}}, {{$outClusterNum}}, EP_{{$endpointID}});				\
	ZB_AF_SIMPLE_DESC_TYPE_VA({{$inClusterNum}}, {{$outClusterNum}}, EP_{{$endpointID}}) simple_desc_##ep_name =	\
	{										\
		{{$endpointID}},									\
		ZB_AF_HA_PROFILE_ID,							\
		ZB_HA_TEMPERATURE_SENSOR_DEVICE_ID,	/*This values are present as initial ones.*/				\
		ZB_HA_DEVICE_VER_TEMPERATURE_SENSOR, /*This values are present as initial ones.*/				\
		0,									\
		{{$inClusterNum}},								\
		{{$outClusterNum}},								\
		{									\
			{{- range $sensor.Clusters}}
			{{.ID.ToZCL}}, \
			{{- end}}
		}									\
	}
{{- end }}

// Define a single endpoint.
// `ep_name` is variable that will be created, and will be later
// used in `ZBOSS_DECLARE_DEVICE_CTX...` macros.
// `cluster_list` is variable created by `ZB_HA_DECLARE_DEVICE_CLUSTER_LIST_...`.
#define ZB_HA_DECLARE_DEVICE_EP(ep_name, ep_id, report_count, cluster_list)				\
	ZB_ZCL_DECLARE_DEVICE_DESC_EP_##ep_id(ep_name);						\
	ZBOSS_DEVICE_DECLARE_REPORTING_CTX(							\
		reporting_info##ep_name,							\
		report_count);					\
	ZB_AF_DECLARE_ENDPOINT_DESC(								\
		ep_name,									\
		ep_id,										\
		ZB_AF_HA_PROFILE_ID,								\
		0,										\
		NULL,										\
		ZB_ZCL_ARRAY_SIZE(cluster_list, zb_zcl_cluster_desc_t),				\
		cluster_list,									\
		(zb_af_simple_desc_1_1_t *)&simple_desc_##ep_name,				\
		report_count, reporting_info##ep_name, 0, NULL)


{{ if gt (len .Device.Sensors) 4}}
// Register endpoints with device ctx
// This macro is only for devices that have >4 endpoints
#define ZBOSS_DECLARE_DEVICE_CTX_{{len .Device.Sensors}}_EP( \
	device_ctx_name, \
	{{- range .Device.Sensors}}
	ep
	{{- end}}
) \
 ZB_AF_START_DECLARE_ENDPOINT_LIST(ep_list_##device_ctx_name) \
 {{- range $i := .Device.Sensors}}
 &ep{{sum $i 1}}_name, \
 {{- end}}
 ZB_AF_FINISH_DECLARE_ENDPOINT_LIST; \
 ZBOSS_DECLARE_DEVICE_CTX(device_ctx_name, ep_list_##device_ctx_name, \
 (ZB_ZCL_ARRAY_SIZE(ep_list_##device_ctx_name, zb_af_endpoint_desc_t*)))
{{ end }}
