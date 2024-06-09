package extenders

import (
	"fmt"
	"log"

	"github.com/ffenix113/zigbee_home/types"
	"github.com/ffenix113/zigbee_home/types/devicetree"
	"github.com/ffenix113/zigbee_home/types/generator"
)

var _ generator.Extender = Button{}
var _ devicetree.Applier = Button{}

type Button struct {
	generator.SimpleExtender

	Instances []types.Pin
}

func NewButtons(instances ...types.Pin) generator.Extender {
	for i := range instances {
		if instances[i].ID == "" {
			log.Fatalf("button %#v must have an id set", instances[i])
		}
	}

	return Button{
		Instances: instances,
	}
}

func (b Button) Template() string {
	return "peripherals/buttons"
}

func (b Button) Includes() []string {
	return []string{"zephyr/drivers/gpio.h"}
}

func (b Button) ApplyOverlay(dt *devicetree.DeviceTree) error {
	for _, instance := range b.Instances {
		ledInstance := devicetree.NewButton(instance)
		if err := ledInstance.AttachSelf(dt); err != nil {
			return fmt.Errorf("attach button: %w", err)
		}
	}

	return nil
}
