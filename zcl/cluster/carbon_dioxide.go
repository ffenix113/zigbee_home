package cluster

var _ Cluster = CarbonDioxide{}

// ZCL 4.5.2
type CarbonDioxide struct {
	MinMeasuredValue int16
	MaxMeasuredValue int16
	Tolerance        uint16
}

func (t CarbonDioxide) ID() ID {
	return ID_CARBON_DIOXIDE
}

func (CarbonDioxide) CAttrType() string {
	return "zb_zcl_measurement_type_single_attrs_t"
}
func (CarbonDioxide) CVarName() string {
	return "carbon_dioxide"
}

func (CarbonDioxide) ReportAttrCount() int {
	return 1
}

func (CarbonDioxide) Side() Side {
	return Server
}
