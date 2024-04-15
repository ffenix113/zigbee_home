package extenders

import (
	"fmt"
	"strings"

	"github.com/ffenix113/zigbee_home/types"
	"github.com/ffenix113/zigbee_home/types/devicetree"
	"github.com/ffenix113/zigbee_home/types/generator"
)

var _ generator.Extender = I2C{}
var _ devicetree.Applier = I2C{}

type I2CInstance struct {
	// ID is a actual label of pre-defined I2C peripheral.
	// For example most SoCs would have IDs something like i2c0, i2c1, ...
	ID       string
	SDA, SCL types.Pin
}

type I2C struct {
	generator.SimpleExtender

	Instances []I2CInstance
}

func NewI2C(instances ...I2CInstance) I2C {
	for i, instance := range instances {
		if len(instance.ID) != 4 || !strings.HasPrefix(instance.ID, "i2c") || (instance.ID[3] < '0' || instance.ID[3] > '9') {
			panic(fmt.Sprintf("i2c instance %d must have `id` format of 'i2c[0-9]'", i))
		}
	}

	return I2C{
		Instances: instances,
	}
}

func (i I2C) ApplyOverlay(dt *devicetree.DeviceTree) error {
	pinctrl := dt.FindSpecificNode(devicetree.SearchByLabel(devicetree.NodeLabelPinctrl))

	for _, instance := range i.Instances {
		// Add pin definitions only if we have some.
		// Otherwise just enable the I2C instance.
		if instance.SDA.PinsDefined() && instance.SCL.PinsDefined() {
			pinctrl.AddNodes(buildI2C(instance.ID, instance)...)
		}

		dt.AddNodes(&devicetree.Node{
			Label:      instance.ID,
			Upsert:     true,
			Properties: []devicetree.Property{devicetree.PropertyStatusEnable},
		})
	}

	return nil
}

func buildI2C(id string, i I2CInstance) []*devicetree.Node {
	return []*devicetree.Node{
		{
			Name:     id + "_default",
			Label:    id + "_default",
			SubNodes: []*devicetree.Node{buildI2CNode(i, false)},
		},
		{
			Name:     id + "_sleep",
			Label:    id + "_sleep",
			SubNodes: []*devicetree.Node{buildI2CNode(i, true)},
		},
	}
}

func buildI2CNode(i I2CInstance, lowPowerEnable bool) *devicetree.Node {
	group1 := &devicetree.Node{
		Name: "group1",
		Properties: []devicetree.Property{
			devicetree.NewProperty("psels",
				devicetree.Array(
					devicetree.NrfPSel("TWIM_SDA", i.SDA.Port.Value(), i.SDA.Pin.Value()),
					devicetree.NrfPSel("TWIM_SCL", i.SCL.Port.Value(), i.SCL.Pin.Value()),
				),
			),
		},
	}

	if lowPowerEnable {
		group1.Properties = append(group1.Properties, devicetree.NewProperty("low-power-enable", nil))
	}

	return group1
}
