package cluster

import "fmt"

var _ Cluster = WaterContent{}

func NewRelativeHumidity(minVal, maxVal uint16) *WaterContent {
	return &WaterContent{
		MinMeasuredValue: minVal,
		MaxMeasuredValue: maxVal,

		ClusterID:   ID_REL_HUMIDITY_MEASUREMENT,
		CVarNameStr: "humidity_measurement",
	}
}

func NewSoilMoisture(minVal, maxVal uint16) *WaterContent {
	return &WaterContent{
		MinMeasuredValue: minVal,
		MaxMeasuredValue: maxVal,

		ClusterID:   ID_SOIL_MOISTURE_MEASUREMENT,
		CVarNameStr: "soil_moisture",
	}
}

// ZCL 4.7.2
type WaterContent struct {
	MinMeasuredValue uint16
	MaxMeasuredValue uint16
	// Tolerance is not supported for humidity in nRF Connect SDK v2.5.0
	// Tolerance uint16

	ClusterID   ID     `yaml:"-"`
	CVarNameStr string `yaml:"-"`
}

func (wc WaterContent) ID() ID {
	return wc.ClusterID
}

func (wc WaterContent) CAttrType() string {
	return fmt.Sprintf("zb_zcl_%s_attrs_t", wc.CVarNameStr)
}
func (wc WaterContent) CVarName() string {
	return wc.CVarNameStr
}

func (WaterContent) ReportAttrCount() int {
	return 1
}

func (WaterContent) Side() Side {
	return Server
}
