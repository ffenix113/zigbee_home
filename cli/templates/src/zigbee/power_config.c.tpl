{{ define "power_config_attr_list" }}

// Define fixed attr list, as normal list does not contain
// battery percentage, while extended is bugged:
// does not provide `batt_num` value to macro.
#define ZB_ZCL_DECLARE_POWER_CONFIG_BATTERY_ATTRIB_LIST(attr_list, voltage, rated_voltage, alarm_mask, voltage_min_threshold, percentage_remaining) \
 ZB_ZCL_START_DECLARE_ATTRIB_LIST_CLUSTER_REVISION(attr_list, ZB_ZCL_POWER_CONFIG) \
 ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_POWER_CONFIG_BATTERY_VOLTAGE_ID(voltage, ), \
 ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_POWER_CONFIG_BATTERY_RATED_VOLTAGE_ID(rated_voltage, ), \
 ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_POWER_CONFIG_BATTERY_ALARM_MASK_ID(alarm_mask, ), \
 ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_POWER_CONFIG_BATTERY_VOLTAGE_MIN_THRESHOLD_ID(voltage_min_threshold, ), \
 ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_POWER_CONFIG_BATTERY_PERCENTAGE_REMAINING_ID(percentage_remaining, ), \
 ZB_ZCL_FINISH_DECLARE_ATTRIB_LIST


ZB_ZCL_DECLARE_POWER_CONFIG_BATTERY_ATTRIB_LIST(
	{{.Cluster.CVarName}}_{{.Endpoint}}_attr_list,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.voltage,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.rated_voltage,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.alarm_mask,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.voltage_min_threshold,
	&dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.percentage_remaining);
{{ end }}

{{ define "power_config_attr_init"}}
	/* IAS Zone */
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.voltage = ZB_ZCL_POWER_CONFIG_BATTERY_VOLTAGE_INVALID;
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.rated_voltage = {{.Cluster.BatteryRatedVoltage | formatHex}};
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.alarm_mask = ZB_ZCL_POWER_CONFIG_BATTERY_ALARM_MASK_DEFAULT_VALUE;
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.voltage_min_threshold = {{.Cluster.BatteryVoltageMinThreshold | formatHex}};
	dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.percentage_remaining = ZB_ZCL_POWER_CONFIG_BATTERY_REMAINING_UNKNOWN;
{{end}}