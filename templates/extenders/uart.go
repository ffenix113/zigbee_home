package extenders

import (
	"fmt"
	"strings"

	"github.com/ffenix113/zigbee_home/types"
	"github.com/ffenix113/zigbee_home/types/appconfig"
	"github.com/ffenix113/zigbee_home/types/devicetree"
	"github.com/ffenix113/zigbee_home/types/generator"
)

var _ appconfig.Provider = UART{}
var _ devicetree.Applier = UART{}

type UARTInstance struct {
	// ID is a actual label of pre-defined UART peripheral.
	// For example most SoCs would have IDs something like uart0, uart1, ...
	ID     string
	TX, RX types.Pin
	Speed  int
}

type UART struct {
	generator.SimpleExtender

	Instances []UARTInstance
}

func NewUART(instances ...UARTInstance) UART {
	for i, instance := range instances {
		if len(instance.ID) != 5 || !strings.HasPrefix(instance.ID, "uart") || (instance.ID[4] < '0' || instance.ID[4] > '9') {
			panic(fmt.Sprintf("uart instance %d(%q) must have `id` format of 'uart[0-9]'", i, instance.ID))
		}
	}

	return UART{
		Instances: instances,
	}
}

func (i UART) AppConfig() []appconfig.ConfigValue {
	return []appconfig.ConfigValue{
		appconfig.CONFIG_UART_CONSOLE.Required(appconfig.Yes),
		appconfig.CONFIG_UART_LINE_CTRL.Required(appconfig.Yes),
	}
}

func (i UART) ApplyOverlay(dt *devicetree.DeviceTree) error {
	pinctrl := dt.FindSpecificNode(devicetree.SearchByLabel(devicetree.NodeLabelPinctrl))

	for _, instance := range i.Instances {
		if instance.Speed <= 0 {
			instance.Speed = 115200
		}

		// Pinctrl

		// Add pin definitions only if we have some.
		// Otherwise just enable the UART instance.
		if instance.RX.PinsDefined() && instance.TX.PinsDefined() {
			pinctrl.AddNodes(buildPinctrlUART(instance.ID, instance)...)
		}

		// Root item
		// Define, even if corrently defined on base board overlay,
		// just to be sure.
		dt.AddNodes(&devicetree.Node{
			Label:  instance.ID,
			Upsert: true,
			Properties: []devicetree.Property{
				devicetree.NewProperty(devicetree.PropertyNameCompatible, devicetree.FromValue("nordic,nrf-uart")),
				devicetree.NewProperty("current-speed", devicetree.FromValue(instance.Speed)),
				devicetree.PropertyStatusEnable,
				devicetree.NewProperty("pinctrl-0", devicetree.Angled(devicetree.Label(fmt.Sprintf("%s_default", instance.ID)))),
				devicetree.NewProperty("pinctrl-1", devicetree.Angled(devicetree.Label(fmt.Sprintf("%s_sleep", instance.ID)))),
				devicetree.NewProperty("pinctrl-names", devicetree.Array(devicetree.Quoted("default"), devicetree.Quoted("sleep"))),
			},
		})
	}

	return nil
}

func buildPinctrlUART(id string, i UARTInstance) []*devicetree.Node {
	return []*devicetree.Node{
		{
			Name:     id + "_default",
			Label:    id + "_default",
			SubNodes: []*devicetree.Node{buildUARTNode(i, false)},
		},
		{
			Name:     id + "_sleep",
			Label:    id + "_sleep",
			SubNodes: []*devicetree.Node{buildUARTNode(i, true)},
		},
	}
}

func buildUARTNode(i UARTInstance, lowPowerEnable bool) *devicetree.Node {
	group1 := &devicetree.Node{
		Name: "group1",
		Properties: []devicetree.Property{
			devicetree.NewProperty("psels",
				devicetree.Array(
					devicetree.NrfPSel("UART_TX", i.TX.Port.Value(), i.TX.Pin.Value()),
					devicetree.NrfPSel("UART_RX", i.RX.Port.Value(), i.RX.Pin.Value()),
				),
			),
		},
	}

	if lowPowerEnable {
		group1.Properties = append(group1.Properties, devicetree.NewProperty("low-power-enable", nil))
	}

	return group1
}
