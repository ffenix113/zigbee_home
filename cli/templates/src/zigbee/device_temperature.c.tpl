{{ define "device_temp_config_attr_list" }}
// ZCL spec 3.4.1.1
#define ZB_ZCL_DEVICE_TEMP_CONFIG_CLUSTER_REVISION_DEFAULT ((zb_uint16_t)0x0001u)

#define ZB_ZCL_ATTR_DEVICE_TEMP_CONFIG_VALUE_UNKNOWN ZB_ZCL_ATTR_PRESSURE_MEASUREMENT_VALUE_UNKNOWN

/*! @brief CurrentTemperature, ZCL spec 3.4.2.2.1 */
#define ZB_ZCL_ATTR_DEVICE_TEMP_CONFIG_CURRENT_TEMPERATURE_ID (0x0000)
// ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_CARBON_DIOXIDE_MIN_VALUE_ID
void zb_zcl_device_temp_config_init_server()
{
  zb_zcl_add_cluster_handlers(ZB_ZCL_CLUSTER_ID_DEVICE_TEMP_CONFIG,
                              ZB_ZCL_CLUSTER_SERVER_ROLE,
                              (zb_zcl_cluster_check_value_t)NULL,
                              (zb_zcl_cluster_write_attr_hook_t)NULL,
                              (zb_zcl_cluster_handler_t)NULL);
}

#define ZB_ZCL_CLUSTER_ID_DEVICE_TEMP_CONFIG_SERVER_ROLE_INIT zb_zcl_device_temp_config_init_server
#define ZB_ZCL_CLUSTER_ID_DEVICE_TEMP_CONFIG_CLIENT_ROLE_INIT (NULL)


typedef void * zb_voidp_t;
#define ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_DEVICE_TEMP_CONFIG_CURRENT_TEMPERATURE_ID(data_ptr) \
  { \
    ZB_ZCL_ATTR_DEVICE_TEMP_CONFIG_CURRENT_TEMPERATURE_ID, \
    ZB_ZCL_ATTR_TYPE_U16, \
    ZB_ZCL_ATTR_ACCESS_READ_ONLY, \
    (zb_voidp_t) data_ptr \
  }

#define ZB_ZCL_DECLARE_DEVICE_TEMP_CONFIG_ATTRIB_LIST(attr_list,                  \
    current_temperature)                                     \
  ZB_ZCL_START_DECLARE_ATTRIB_LIST_CLUSTER_REVISION(attr_list, ZB_ZCL_DEVICE_TEMP_CONFIG) \
  ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_DEVICE_TEMP_CONFIG_CURRENT_TEMPERATURE_ID, (current_temperature))          \
  ZB_ZCL_FINISH_DECLARE_ATTRIB_LIST


ZB_ZCL_DECLARE_DEVICE_TEMP_CONFIG_ATTRIB_LIST(
	{{.Cluster.CVarName}}_{{.Endpoint}}_attr_list,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.current_temperature
	);
{{ end }}

{{ define "device_temp_config_attr_init"}}
	/* Device temperature */
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.current_temperature = ZB_ZCL_ATTR_DEVICE_TEMP_CONFIG_VALUE_UNKNOWN;
{{end}}