package cluster

type OnOff struct {
	PinLabel string
}

func (o OnOff) ID() ID {
	return ID_ON_OFF
}

func (OnOff) CAttrType() string {
	return "zb_zcl_on_off_attrs_t"
}
func (OnOff) CVarName() string {
	return "on_off"
}

func (OnOff) ReportAttrCount() int {
	return 1
}

func (OnOff) Side() Side {
	// This only allows to be controlled from the client,
	// i.e. home assistant.
	// But it cannot be controlled from device for now.
	return Server
}
