{{ define "identify_attr_list" }}
ZB_ZCL_DECLARE_IDENTIFY_CLIENT_ATTRIB_LIST(
	{{.Cluster.CVarName}}_client_attr_list);

ZB_ZCL_DECLARE_IDENTIFY_SERVER_ATTRIB_LIST(
	{{.Cluster.CVarName}}_server_attr_list,
	&dev_ctx.{{.Cluster.CVarName}}_attr.identify_time);
{{ end }}