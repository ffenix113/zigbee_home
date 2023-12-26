package cluster

var _ Cluster = Identify{}

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

func (Identify) Reports() bool {
	return true
}
