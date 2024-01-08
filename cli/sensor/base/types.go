package base

import "strings"

type I2CConnection struct {
	ID   string
	Addr string
}

func (c I2CConnection) UnitAddress() string {
	if strings.HasPrefix(c.Addr, "0x") {
		return c.Addr[2:]
	}

	return c.Addr
}

func (c I2CConnection) Reg() string {
	if !strings.HasPrefix(c.Addr, "0x") {
		return "0x" + c.Addr
	}

	return c.Addr
}
