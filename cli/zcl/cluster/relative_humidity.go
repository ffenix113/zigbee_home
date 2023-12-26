package cluster

var _ Cluster = RelativeHumidity{}

type RelativeHumidity struct {
	MinMeasuredValue int16
	MaxMeasuredValue int16
	// Tolerance is not supported for humidity in nRF Connect SDK v2.5.0
	// Tolerance uint16
}

func (t RelativeHumidity) ID() ID {
	return ID_REL_HUMIDITY_MEASUREMENT
}

func (RelativeHumidity) CAttrType() string {
	return "zb_zcl_humidity_measurement_attrs_t"
}
func (RelativeHumidity) CVarName() string {
	return "humidity_measurement"
}

func (RelativeHumidity) Reports() bool {
	return true
}
