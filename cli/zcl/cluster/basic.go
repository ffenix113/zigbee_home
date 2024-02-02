package cluster

var _ Cluster = Basic{}

// ZCL 3.5
type Basic struct{}

func (t Basic) ID() ID {
	return ID_BASIC
}

func (Basic) CAttrType() string {
	return "zb_zcl_basic_attrs_t"
}
func (Basic) CVarName() string {
	return "basic"
}

func (Basic) ReportAttrCount() int {
	return 0
}

func (Basic) Side() Side {
	return Server
}
