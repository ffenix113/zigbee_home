{{ define "carbon_dioxide_defines"}}
// ZCL spec 4.13.1.3
#define ZB_ZCL_CLUSTER_ID_CARBON_DIOXIDE (0x040d)

// ZCL spec 4.13.1.1
#define ZB_ZCL_CARBON_DIOXIDE_CLUSTER_REVISION_DEFAULT ((zb_uint16_t)0x0002u)

#define ZB_ZCL_ATTR_CARBON_DIOXIDE_VALUE_UNKNOWN ZB_ZCL_ATTR_PRESSURE_MEASUREMENT_VALUE_UNKNOWN

#define ZCL_CARBON_DIOXIDE_MEASURED_VALUE_MULTIPLIER 0.000001

/*! @brief CurrentTemperature, ZCL spec 3.4.2.2.1 */
#define ZB_ZCL_ATTR_CARBON_DIOXIDE_VALUE_ID (0x0000)
#define ZB_ZCL_ATTR_CARBON_DIOXIDE_MIN_VALUE_ID (0x0001)
#define ZB_ZCL_ATTR_CARBON_DIOXIDE_MAX_VALUE_ID (0x0002)
#define ZB_ZCL_ATTR_CARBON_DIOXIDE_TOLERANCE_ID (0x0003)
{{end}}

{{ define "carbon_dioxide_attr_types"}}
void zb_zcl_carbon_dioxide_init_server()
{
  zb_zcl_add_cluster_handlers(ZB_ZCL_CLUSTER_ID_CARBON_DIOXIDE,
                              ZB_ZCL_CLUSTER_SERVER_ROLE,
                              (zb_zcl_cluster_check_value_t)NULL,
                              (zb_zcl_cluster_write_attr_hook_t)NULL,
                              (zb_zcl_cluster_handler_t)NULL);
}

#define ZB_ZCL_CLUSTER_ID_CARBON_DIOXIDE_SERVER_ROLE_INIT zb_zcl_carbon_dioxide_init_server
#define ZB_ZCL_CLUSTER_ID_CARBON_DIOXIDE_CLIENT_ROLE_INIT (NULL)


typedef void * zb_voidp_t;
#define ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_CARBON_DIOXIDE_VALUE_ID(data_ptr) \
  { \
    ZB_ZCL_ATTR_CARBON_DIOXIDE_VALUE_ID, \
    ZB_ZCL_ATTR_TYPE_SINGLE, \
    ZB_ZCL_ATTR_ACCESS_READ_ONLY | ZB_ZCL_ATTR_ACCESS_REPORTING, \
    {{ if not ncsVersionIs_2_5 -}}
    (ZB_UINT16_MAX), \
    {{ end -}}
    (zb_voidp_t) data_ptr \
  }

#define ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_CARBON_DIOXIDE_MIN_VALUE_ID(data_ptr) \
  { \
    ZB_ZCL_ATTR_CARBON_DIOXIDE_MIN_VALUE_ID, \
    ZB_ZCL_ATTR_TYPE_SINGLE, \
    ZB_ZCL_ATTR_ACCESS_READ_ONLY, \
    {{ if not ncsVersionIs_2_5 -}}
    (ZB_UINT16_MAX), \
    {{ end -}}
    (zb_voidp_t) data_ptr \
  }

#define ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_CARBON_DIOXIDE_MAX_VALUE_ID(data_ptr) \
  { \
    ZB_ZCL_ATTR_CARBON_DIOXIDE_MAX_VALUE_ID, \
    ZB_ZCL_ATTR_TYPE_SINGLE, \
    ZB_ZCL_ATTR_ACCESS_READ_ONLY, \
    {{ if not ncsVersionIs_2_5 -}}
    (ZB_UINT16_MAX), \
    {{ end -}}
    (zb_voidp_t) data_ptr \
  }

#define ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_CARBON_DIOXIDE_TOLERANCE_ID(data_ptr) \
  { \
    ZB_ZCL_ATTR_CARBON_DIOXIDE_TOLERANCE_ID, \
    ZB_ZCL_ATTR_TYPE_SINGLE, \
    ZB_ZCL_ATTR_ACCESS_READ_ONLY, \
    {{ if not ncsVersionIs_2_5 -}}
    (ZB_UINT16_MAX), \
    {{ end -}}
    (zb_voidp_t) data_ptr \
  }

#define ZB_ZCL_DECLARE_CARBON_DIOXIDE_ATTRIB_LIST(attr_list,                  \
    value, min_value, max_value, tolerance)                                     \
  ZB_ZCL_START_DECLARE_ATTRIB_LIST_CLUSTER_REVISION(attr_list, ZB_ZCL_CARBON_DIOXIDE) \
  ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_CARBON_DIOXIDE_VALUE_ID, (value))          \
  ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_CARBON_DIOXIDE_MIN_VALUE_ID, (min_value))          \
  ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_CARBON_DIOXIDE_MAX_VALUE_ID, (max_value))          \
  ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_CARBON_DIOXIDE_TOLERANCE_ID, (tolerance))          \
  ZB_ZCL_FINISH_DECLARE_ATTRIB_LIST

{{end}}

{{ define "carbon_dioxide_attr_list" }}
ZB_ZCL_DECLARE_CARBON_DIOXIDE_ATTRIB_LIST(
	{{.Cluster.CVarName}}_{{.Endpoint}}_attr_list,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.measure_value,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.min_measure_value,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.max_measure_value,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.tolerance
	);
{{ end }}

{{ define "carbon_dioxide_attr_init"}}
	/* Carbon Dioxide */
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.measure_value = ZB_ZCL_ATTR_TEMP_MEASUREMENT_VALUE_UNKNOWN;
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.min_measure_value = ({{.Cluster.MinMeasuredValue}} * ZCL_CARBON_DIOXIDE_MEASURED_VALUE_MULTIPLIER);
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.max_measure_value = ({{.Cluster.MaxMeasuredValue}} * ZCL_CARBON_DIOXIDE_MEASURED_VALUE_MULTIPLIER);
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.tolerance = ({{.Cluster.Tolerance}} * ZCL_CARBON_DIOXIDE_MEASURED_VALUE_MULTIPLIER);
{{end}}