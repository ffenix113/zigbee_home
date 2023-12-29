package extenders

import (
	"github.com/ffenix113/zigbee_home/cli/types/devicetree"
	"github.com/ffenix113/zigbee_home/cli/types/generator"
)

var _ generator.Extender = I2C{}
var _ devicetree.Applier = I2C{}

type I2C struct {
	generator.SimpleExtender `yaml:"-"`

	Port     int
	SDA, SCL int
}

func NewI2C(port, sda, scl int) I2C {
	return I2C{
		Port: port,
		SDA:  sda,
		SCL:  scl,
	}
}

func (i I2C) ApplyOverlay(dt *devicetree.DeviceTree) error {
	group1 := &devicetree.Node{
		Name: "group1",
		Properties: []devicetree.Property{
			devicetree.NewProperty("psels",
				devicetree.Array(
					devicetree.NrfPSel("TWIM_SDA", i.Port, i.SDA),
					devicetree.NrfPSel("TWIM_SCL", i.Port, i.SCL),
				),
			),
		},
	}
	group1Sleep := &devicetree.Node{
		Name: "group1",
		Properties: []devicetree.Property{
			devicetree.NewProperty("psels",
				devicetree.Array(
					devicetree.NrfPSel("TWIM_SDA", i.Port, i.SDA),
					devicetree.NrfPSel("TWIM_SCL", i.Port, i.SCL),
				),
			),
			devicetree.NewProperty("low-power-enable", nil),
		},
	}

	pinctrl := dt.FindSpecificNode(devicetree.SearchByLabel("pinctrl"))
	pinctrl.AddNodes(
		&devicetree.Node{
			Name:     "i2c1_default",
			Label:    "i2c1_default",
			SubNodes: []*devicetree.Node{group1},
		},
		&devicetree.Node{
			Name:     "i2c1_sleep",
			Label:    "i2c1_sleep",
			SubNodes: []*devicetree.Node{group1Sleep},
		},
	)

	return nil
}
