{{ define "ias_zone_attr_list" }}
// For some reason non-extended version did not work
// for me, but extended does work fine.
ZB_ZCL_DECLARE_IAS_ZONE_ATTRIB_LIST_EXT(
	{{.Cluster.CVarName}}_{{.Endpoint}}_attr_list,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.zone_state,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.zone_type,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.zone_status,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.number_of_zone_sens_levels_supported,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.current_zone_sens_level,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.ias_cie_address,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.zone_id,
    &dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.cie_short_addr,
    &dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.cie_ep
	);
{{end}}

{{ define "ias_zone_attr_init"}}
	/* IAS Zone */
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.zone_state = ZB_ZCL_IAS_ZONE_ZONESTATE_DEF_VALUE;
    dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.zone_type = {{.Cluster.ZoneType}};
	// Set status to include Restore, to specify that 
	// contact sensor does know when it is closed & open.
    dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.zone_status = ZB_ZCL_IAS_ZONE_ZONE_STATUS_RESTORE;
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.number_of_zone_sens_levels_supported = ZB_ZCL_IAS_ZONE_NUMBER_OF_ZONE_SENSITIVITY_LEVELS_SUPPORTED_DEFAULT_VALUE;
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.current_zone_sens_level = ZB_ZCL_IAS_ZONE_CURRENT_ZONE_SENSITIVITY_LEVEL_DEFAULT_VALUE;
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.zone_id = ZB_ZCL_IAS_ZONEID_ID_DEF_VALUE;
{{end}}