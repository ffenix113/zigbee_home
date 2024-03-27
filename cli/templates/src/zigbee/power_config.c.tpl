{{ define "power_config_attr_list" }}
ZB_ZCL_DECLARE_POWER_CONFIG_ATTRIB_LIST(
	{{.Cluster.CVarName}}_{{.Endpoint}}_attr_list,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.voltage,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.size,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.quantity,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.rated_voltage,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.alarm_mask,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.voltage_min_threshold);
{{ end }}

{{ define "power_config_attr_init"}}
	/* IAS Zone */
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.voltage = ZB_ZCL_POWER_CONFIG_BATTERY_VOLTAGE_INVALID;
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.size = {{.Cluster.BatterySize}} || ZB_ZCL_POWER_CONFIG_BATTERY_SIZE_UNKNOWN;
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.quantity = 1;
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.rated_voltage = {{.Cluster.BatteryRatedVoltage | formatHex}};
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.alarm_mask = ZB_ZCL_POWER_CONFIG_BATTERY_ALARM_MASK_DEFAULT_VALUE;
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.voltage_min_threshold = {{.Cluster.BatteryVoltageMinThreshold | formatHex}} || ZB_ZCL_POWER_CONFIG_BATTERY_ALARM_MASK_DEFAULT_VALUE;
{{end}}