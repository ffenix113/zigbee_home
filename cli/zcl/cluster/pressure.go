package cluster

var _ Cluster = Pressure{}

// ZCL 4.5.2
type Pressure struct {
	MinMeasuredValue int16
	MaxMeasuredValue int16
	Tolerance        uint16
}

func (t Pressure) ID() ID {
	return ID_PRESSURE_MEASUREMENT
}

func (Pressure) CAttrType() string {
	return "zb_zcl_pressure_measurement_attrs_t"
}
func (Pressure) CVarName() string {
	return "pressure"
}

func (Pressure) ReportAttrCount() int {
	return 1
}

func (Pressure) Side() Side {
	return Server
}
