package cluster

var _ Cluster = Pressure{}

type Pressure struct{}

func (t Pressure) ID() ID {
	return ID_PRESSURE_MEASUREMENT
}

func (Pressure) CAttrType() string {
	return "zb_zcl_pressure_measurement_attrs_t"
}
func (Pressure) CVarName() string {
	return "pressure_measurement"
}

func (Pressure) Reports() bool {
	return true
}
