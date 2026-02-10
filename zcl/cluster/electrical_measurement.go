package cluster

var _ Cluster = ElectricalMeasurement{}

// ZCL 4.9
type ElectricalMeasurement struct{}

func (t ElectricalMeasurement) ID() ID {
	return ID_ELECTRICAL_MEASUREMENT
}

func (ElectricalMeasurement) CAttrType() string {
	return "zb_zcl_electrical_measurement_attrs_t"
}
func (ElectricalMeasurement) CVarName() string {
	return "electrical_measurement"
}

func (ElectricalMeasurement) ReportAttrCount() int {
	// Quite a lot, but this is so Z2M will not error out when trying to configure it.
	// https://github.com/Koenkk/zigbee-herdsman-converters/issues/11367
	return 16
}

func (ElectricalMeasurement) Side() Side {
	// Can be added as client as well,
	// when supported by templates.
	return Server
}

func (ElectricalMeasurement) CPPArgs() []any {
	return nil
}
