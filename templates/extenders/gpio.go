package extenders

import (
	"github.com/ffenix113/zigbee_home/types/generator"
)

var _ generator.Extender = GPIO{}

type GPIO struct {
	generator.SimpleExtender `yaml:"-"`
}

func (GPIO) Includes() []string {
	return []string{"zephyr/drivers/gpio.h"}
}
