{{ define "water_content_defines"}}
// ZCL spec 4.7.1.3
#define ZB_ZCL_CLUSTER_ID_SOIL_MOISTURE (0x0408)

// ZCL spec 4.7.1.1
#define ZB_ZCL_SOIL_MOISTURE_CLUSTER_REVISION_DEFAULT ((zb_uint16_t)0x0002u)

#define ZB_ZCL_ATTR_SOIL_MOISTURE_VALUE_UNKNOWN (0xffff)

{{ end }}

{{ define "water_content_attr_list" }}

void zb_zcl_soil_moisture_init_server()
{
  zb_zcl_add_cluster_handlers(ZB_ZCL_CLUSTER_ID_SOIL_MOISTURE,
                              ZB_ZCL_CLUSTER_SERVER_ROLE,
                              (zb_zcl_cluster_check_value_t)NULL,
                              (zb_zcl_cluster_write_attr_hook_t)NULL,
                              (zb_zcl_cluster_handler_t)NULL);
}

#define ZB_ZCL_CLUSTER_ID_SOIL_MOISTURE_SERVER_ROLE_INIT zb_zcl_soil_moisture_init_server
#define ZB_ZCL_CLUSTER_ID_SOIL_MOISTURE_CLIENT_ROLE_INIT (NULL)

ZB_ZCL_DECLARE_REL_HUMIDITY_MEASUREMENT_ATTRIB_LIST(
	{{.Cluster.CVarName}}_{{.Endpoint}}_attr_list,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.measure_value,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.min_measure_value,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.max_measure_value
	);
{{ end }}

{{define "water_content_attr_init"}}
	/* Water content, {{.Cluster.CVarName}} */
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.measure_value = ZB_ZCL_ATTR_REL_HUMIDITY_MEASUREMENT_VALUE_UNKNOWN;
	// Range is 0 - 100%
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.min_measure_value = 0;
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.max_measure_value = 100 * 100;	
	/* Humidity measurements tolerance is not supported at the moment */
{{end}}