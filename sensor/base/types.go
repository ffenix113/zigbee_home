package base

import (
	"fmt"
)

type I2CConnection struct {
	ID   string
	Addr uint8
}

func (c I2CConnection) UnitAddress() string {
	return fmt.Sprintf("%x", c.Addr)
}

func (c I2CConnection) Reg() string {
	return fmt.Sprintf("%#x", c.Addr)
}
