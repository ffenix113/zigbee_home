{{ define "level_control_attr_list" }}
// Here we don't need to have exntender, as for now usage is limited
// to only sending command(s) to server.
ZB_ZCL_DECLARE_LEVEL_CONTROL_ATTRIB_LIST(
	{{.Cluster.CVarName}}_{{.Endpoint}}_attr_list,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.current_level,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.remaining_time
);
{{end}}

{{ define "level_control_attr_init"}}
	/* Level control */
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.current_level = ZB_ZCL_LEVEL_CONTROL_CURRENT_LEVEL_DEFAULT_VALUE;
    dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.remaining_time = ZB_ZCL_LEVEL_CONTROL_REMAINING_TIME_DEFAULT_VALUE;
{{end}}