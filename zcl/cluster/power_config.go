package cluster

type PowerConfiguration struct {
	BatteryRatedVoltage        uint16 `yaml:"battery_rated_voltage"`
	BatteryVoltageMinThreshold uint16 `yaml:"battery_voltage_min_threshold"`
}

func (o PowerConfiguration) ID() ID {
	return ID_POWER_CONFIG
}

func (PowerConfiguration) CAttrType() string {
	return "zb_zcl_power_config_attrs_t"
}
func (PowerConfiguration) CVarName() string {
	return "power_config"
}

func (PowerConfiguration) ReportAttrCount() int {
	return 1
}

func (PowerConfiguration) Side() Side {
	return Server
}
