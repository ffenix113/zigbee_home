{{ define "on_off_attr_list" }}
ZB_ZCL_DECLARE_ON_OFF_ATTRIB_LIST(
	{{.Cluster.CVarName}}_{{.Endpoint}}_attr_list,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.on_off);
{{ end }}