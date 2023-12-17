{{ define "define_on_off_attr_list" }}
ZB_ZCL_DECLARE_ON_OFF_ATTRIB_LIST(
	{{.CVarName}}_attr_list,
	&dev_ctx.{{.CVarName}}_attrs.on_off);
{{ end }}