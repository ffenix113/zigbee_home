package extenders

import (
	"github.com/ffenix113/zigbee_home/cli/types/appconfig"
	"github.com/ffenix113/zigbee_home/cli/types/devicetree"
	"github.com/ffenix113/zigbee_home/cli/types/generator"
)

var _ appconfig.Provider = USBUARTLog{}
var _ devicetree.Applier = USBUARTLog{}

type USBUARTLog struct {
	generator.SimpleExtender
}

func NewUSBUARTLog() generator.Extender {
	return USBUARTLog{
		SimpleExtender: generator.SimpleExtender{
			IncludeHeaders: []string{
				"zephyr/drivers/uart.h",
				"zephyr/usb/usb_device.h",
				"zephyr/usb/usbd.h",
				"zephyr/sys/printk.h",
			},
		},
	}
}

func (USBUARTLog) AppConfig() []appconfig.ConfigValue {
	return []appconfig.ConfigValue{
		// Logging setup
		appconfig.CONFIG_LOG.Required(appconfig.Yes),
		appconfig.CONFIG_CONSOLE.Required(appconfig.Yes),
		appconfig.CONFIG_SERIAL.Required(appconfig.Yes),
		appconfig.CONFIG_UART_CONSOLE.Required(appconfig.Yes),
		appconfig.CONFIG_UART_LINE_CTRL.Required(appconfig.Yes),

		// USB
		appconfig.CONFIG_USB_DEVICE_STACK.Required(appconfig.Yes),
	}
}

// ApplyOverlay implements devicetree.Applier.
func (USBUARTLog) ApplyOverlay(dt *devicetree.DeviceTree) error {
	udc := dt.FindSpecificNode(devicetree.SearchByLabel("zephyr_udc0"))
	if udc == nil {
		dt = dt.AddNodes(&devicetree.Node{
			Label:  "zephyr_udc0",
			Upsert: true,
		})
		udc = dt.FindSpecificNode(devicetree.SearchByLabel("zephyr_udc0"))
	}

	udc.AddNodes(&devicetree.Node{
		Name:  "cdc_acm_uart0",
		Label: "cdc_acm_uart0",
		Properties: []devicetree.Property{
			devicetree.NewProperty(devicetree.PropertyNameCompatible, devicetree.FromValue("zephyr,cdc-acm-uart")),
		},
	})

	chosen := dt.FindSpecificNode(
		devicetree.SearchByName(devicetree.NodeNameRoot),
		devicetree.SearchByName(devicetree.NodeNameChosen))

	chosen.AddProperties(
		devicetree.NewProperty("zephyr,console", devicetree.Label("cdc_acm_uart0")),
	)

	return nil
}
