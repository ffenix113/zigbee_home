package cluster

var _ Cluster = Temperature{}

type Temperature struct{}

func (t Temperature) ID() ID {
	return ID_TEMP_MEASUREMENT
}

func (Temperature) CAttrType() string {
	return "zb_zcl_temp_measurement_attrs_t"
}
func (Temperature) CVarName() string {
	return "temp_measurement"
}

func (Temperature) Reports() bool {
	return true
}
