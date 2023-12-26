{{ define "define_identify_attr_list" }}
ZB_ZCL_DECLARE_IDENTIFY_CLIENT_ATTRIB_LIST(
	{{.CVarName}}_client_attr_list);

ZB_ZCL_DECLARE_IDENTIFY_SERVER_ATTRIB_LIST(
	{{.CVarName}}_server_attr_list,
	&dev_ctx.{{.CVarName}}_attr.identify_time);
{{ end }}