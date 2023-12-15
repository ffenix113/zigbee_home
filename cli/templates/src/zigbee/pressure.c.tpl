{{ define "define_temp_attr_list" }}
ZB_ZCL_DECLARE_PRESSURE_MEASUREMENT_ATTRIB_LIST(
	{{clusterToCVar .}}_attr_list,
	&dev_ctx.pres_attrs.measure_value,
	&dev_ctx.pres_attrs.min_measure_value,
	&dev_ctx.pres_attrs.max_measure_value,
	&dev_ctx.pres_attrs.tolerance
	);
{{ end }}