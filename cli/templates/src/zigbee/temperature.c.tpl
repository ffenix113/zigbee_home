{{ define "define_temperature_attr_list" }}
ZB_ZCL_DECLARE_TEMP_MEASUREMENT_ATTRIB_LIST(
	{{.Cluster.CVarName}}_{{.Endpoint}}_attr_list,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.measure_value,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.min_measure_value,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.max_measure_value,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.tolerance
	);
{{ end }}