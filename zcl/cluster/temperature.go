package cluster

var _ Cluster = Temperature{}

// ZCL 4.4.2
type Temperature struct {
	MinMeasuredValue int16
	MaxMeasuredValue int16
	Tolerance        uint16
}

func (t Temperature) ID() ID {
	return ID_TEMP_MEASUREMENT
}

func (Temperature) CAttrType() string {
	return "zb_zcl_temp_measurement_attrs_t"
}
func (Temperature) CVarName() string {
	return "temperature"
}

func (Temperature) ReportAttrCount() int {
	return 1
}

func (Temperature) Side() Side {
	return Server
}
