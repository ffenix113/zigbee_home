package cluster

var _ Cluster = IASZone{}

type IasZoneType string

const (
	IasZoneContact IasZoneType = "contact"
)

func (c IasZoneType) String() string {
	var suffix string

	switch c {
	case IasZoneContact:
		suffix = "CONTACT_SWITCH"
	default:
		// For now this would work
		suffix = "STANDARD_CIE"
	}

	return "ZB_ZCL_IAS_ZONE_ZONETYPE_" + suffix
}

// ZCL 4.5.2
//
// Currently only supports contact zone type.
type IASZone struct {
	ZoneType IasZoneType
}

func (z IASZone) ID() ID {
	return ID_IAS_ZONE
}

func (IASZone) CAttrType() string {
	return "zb_zcl_ias_zone_attrs_t"
}
func (IASZone) CVarName() string {
	return "ias_zone"
}

func (IASZone) ReportAttrCount() int {
	return 0
}

func (z IASZone) Side() Side {
	return Server
}
