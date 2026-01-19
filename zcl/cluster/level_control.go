package cluster

var _ Cluster = LevelControl{}

// ZCL 3.10.2
type LevelControl struct {
}

func (z LevelControl) ID() ID {
	return ID_LEVEL_CONTROL
}

func (LevelControl) CAttrType() string {
	return "zb_zcl_level_control_attrs_t"
}
func (LevelControl) CVarName() string {
	return "level_control"
}

func (LevelControl) ReportAttrCount() int {
	return 0
}

func (z LevelControl) Side() Side {
	return Client
}
