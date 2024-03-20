package extenders

import (
	"fmt"

	"github.com/ffenix113/zigbee_home/cli/types"
	"github.com/ffenix113/zigbee_home/cli/types/devicetree"
	"github.com/ffenix113/zigbee_home/cli/types/generator"
)

var _ generator.Extender = LED{}
var _ devicetree.Applier = LED{}

type LED struct {
	generator.SimpleExtender

	Instances []types.Pin
}

func NewLEDs(instances ...types.Pin) generator.Extender {
	return LED{
		Instances: instances,
	}
}

func (l LED) Template() string {
	return "peripherals/leds"
}

func (l LED) Includes() []string {
	return []string{"zephyr/drivers/gpio.h"}
}

func (l LED) ApplyOverlay(dt *devicetree.DeviceTree) error {
	for _, instance := range l.Instances {
		ledInstance := devicetree.NewLED(instance)
		if err := ledInstance.AttachSelf(dt); err != nil {
			return fmt.Errorf("attach led: %w", err)
		}
	}

	return nil
}
