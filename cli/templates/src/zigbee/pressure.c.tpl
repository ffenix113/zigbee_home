{{ define "pressure_attr_list" }}
ZB_ZCL_DECLARE_PRESSURE_MEASUREMENT_ATTRIB_LIST(
	{{.Cluster.CVarName}}_{{.Endpoint}}_attr_list,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.measure_value,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.min_measure_value,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.max_measure_value,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.tolerance
	);
{{ end }}

{{ define "pressure_attr_init"}}
	/* Pressure */
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.measure_value = ZB_ZCL_ATTR_PRESSURE_MEASUREMENT_VALUE_UNKNOWN;
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.min_measure_value = ({{.Cluster.MinMeasuredValue}} * ZCL_PRESSURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER);
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.max_measure_value = ({{.Cluster.MaxMeasuredValue}} * ZCL_PRESSURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER);
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.tolerance = ({{.Cluster.Tolerance}} * ZCL_PRESSURE_MEASUREMENT_MEASURED_VALUE_MULTIPLIER);
{{end}}