{{ define "define_humidity_attr_list" }}
ZB_ZCL_DECLARE_REL_HUMIDITY_MEASUREMENT_ATTRIB_LIST(
	{{.CVarName}}_attr_list,
	&dev_ctx.{{.CVarName}}_attrs.measure_value,
	&dev_ctx.{{.CVarName}}_attrs.min_measure_value,
	&dev_ctx.{{.CVarName}}_attrs.max_measure_value
	);
{{ end }}