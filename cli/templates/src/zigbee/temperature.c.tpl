{{ define "define_temp_attr_list" }}
ZB_ZCL_DECLARE_TEMP_MEASUREMENT_ATTRIB_LIST(
	{{clusterToCVar .}}_attr_list,
	&dev_ctx.temp_attrs.measure_value,
	&dev_ctx.temp_attrs.min_measure_value,
	&dev_ctx.temp_attrs.max_measure_value,
	&dev_ctx.temp_attrs.tolerance
	);
{{ end }}