package cluster

var _ Cluster = DeviceTemperature{}

// ZCL 3.4.2
type DeviceTemperature struct{}

func (t DeviceTemperature) ID() ID {
	return ID_DEVICE_TEMP_CONFIG
}

func (DeviceTemperature) CAttrType() string {
	return "zb_zcl_device_temperature_config_attrs_t"
}
func (DeviceTemperature) CVarName() string {
	return "device_temperature"
}

func (DeviceTemperature) ReportAttrCount() int {
	return 0
}

func (DeviceTemperature) Side() Side {
	return Server
}
