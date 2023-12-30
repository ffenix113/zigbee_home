package cluster

var _ Cluster = Identify{}

// ZCL 3.5
type Identify struct{}

func (t Identify) ID() ID {
	return ID_IDENTIFY
}

func (Identify) CAttrType() string {
	return "zb_zcl_identify_attrs_t"
}
func (Identify) CVarName() string {
	return "identify"
}

func (Identify) ReportAttrCount() int {
	return 0
}

func (Identify) Side() Side {
	// Can be added as client as well,
	// when supported by templates.
	return Server
}
