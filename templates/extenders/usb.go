package extenders

import (
	"github.com/ffenix113/zigbee_home/types/appconfig"
	"github.com/ffenix113/zigbee_home/types/devicetree"
	"github.com/ffenix113/zigbee_home/types/generator"
)

var _ generator.Extender = UART{}
var _ devicetree.Applier = UART{}

// NewUSBUART will add necessary configuration to enable
// USB UART output.
//
// https://docs.zephyrproject.org/latest/connectivity/usb/device/usb_device.html#console-over-cdc-acm-uart
func NewUSBUART() generator.Extender {
	return generator.SimpleExtender{
		IncludeHeaders: []string{
			"zephyr/usb/usb_device.h",
			"zephyr/usb/usbd.h",
		},
		Config: []appconfig.ConfigValue{
			appconfig.CONFIG_USB_DEVICE_STACK.Required(appconfig.Yes),
			appconfig.CONFIG_USB_DEVICE_INITIALIZE_AT_BOOT.Required(appconfig.No),
		},
		OverlayFn: applyOverlay,
	}
}

func applyOverlay(dt *devicetree.DeviceTree) error {
	dt.AddNodes(&devicetree.Node{
		Label:  "zephyr_udc0",
		Upsert: true,
		SubNodes: []*devicetree.Node{
			{
				Name:  "cdc_acm_uart0",
				Label: "cdc_acm_uart0",
				Properties: []devicetree.Property{
					devicetree.NewProperty(devicetree.PropertyNameCompatible, devicetree.Quoted("zephyr,cdc-acm-uart")),
				},
			},
		},
	})

	return nil
}
