{{ define "define_on_off_attr_list" }}
ZB_ZCL_DECLARE_ON_OFF_ATTRIB_LIST(
	{{clusterToCVar .}}_attr_list,
	&dev_ctx.on_off_attr.on_off);
{{ end }}