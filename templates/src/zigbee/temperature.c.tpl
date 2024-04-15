{{ define "temperature_attr_list" }}
ZB_ZCL_DECLARE_TEMP_MEASUREMENT_ATTRIB_LIST(
	{{.Cluster.CVarName}}_{{.Endpoint}}_attr_list,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.measure_value,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.min_measure_value,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.max_measure_value,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.tolerance
	);
{{ end }}

{{ define "temperature_attr_init"}}
	/* Temperature */
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.measure_value = ZB_ZCL_ATTR_TEMP_MEASUREMENT_VALUE_UNKNOWN;
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.min_measure_value = ({{.Cluster.MinMeasuredValue}} * ZCL_TEMPERATURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER);
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.max_measure_value = ({{.Cluster.MaxMeasuredValue}} * ZCL_TEMPERATURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER);
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.tolerance = ({{.Cluster.Tolerance}} * ZCL_TEMPERATURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER);
{{ end}}