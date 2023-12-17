{{ define "define_temperature_attr_list" }}
ZB_ZCL_DECLARE_TEMP_MEASUREMENT_ATTRIB_LIST(
	{{.CVarName}}_attr_list,
	&dev_ctx.{{.CVarName}}_attrs.measure_value,
	&dev_ctx.{{.CVarName}}_attrs.min_measure_value,
	&dev_ctx.{{.CVarName}}_attrs.max_measure_value,
	&dev_ctx.{{.CVarName}}_attrs.tolerance
	);
{{ end }}