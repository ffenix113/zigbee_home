{{ define "humidity_attr_list" }}
ZB_ZCL_DECLARE_REL_HUMIDITY_MEASUREMENT_ATTRIB_LIST(
	{{.Cluster.CVarName}}_{{.Endpoint}}_attr_list,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.measure_value,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.min_measure_value,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.max_measure_value
	);
{{ end }}

{{define "humidity_attr_init"}}
	/* Humidity */
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.measure_value = ZB_ZCL_ATTR_REL_HUMIDITY_MEASUREMENT_VALUE_UNKNOWN;
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.min_measure_value = ({{.Cluster.MinMeasuredValue}} * ZCL_HUMIDITY_MEASUREMENT_MEASURED_VALUE_MULTIPLIER);
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.max_measure_value = ({{.Cluster.MaxMeasuredValue}} * ZCL_HUMIDITY_MEASUREMENT_MEASURED_VALUE_MULTIPLIER);	
	/* Humidity measurements tolerance is not supported at the moment */
{{end}}