package base

import (
	"github.com/ffenix113/zigbee_home/templates/extenders"
	"github.com/ffenix113/zigbee_home/types"
	"github.com/ffenix113/zigbee_home/types/appconfig"
	"github.com/ffenix113/zigbee_home/types/devicetree"
	"github.com/ffenix113/zigbee_home/types/generator"
	"github.com/ffenix113/zigbee_home/zcl/cluster"
)

type OnOff struct {
	*Base `yaml:",inline"`
	Pin   types.Pin
	// Button will be a on/off physical switch somewhere.
	// For example a button on the board.
	// This will also change the state of this OnOff.
	//
	// Take an example of a power socket:
	// you may want to control it with Home Assistant,
	// but you also sure want to control it from the device itself.
	ControlButton string `yaml:"control_button"`
}

func (*OnOff) String() string {
	return "On/Off"
}

func (*OnOff) Template() string {
	return "sensors/on_off"
}

func (*OnOff) CPPComponentType() string {
	return "OnOff"
}

func (o *OnOff) Clusters() cluster.Clusters {
	return []cluster.Cluster{
		cluster.OnOff{PinLabel: o.Pin.Label()},
	}
}

func (*OnOff) WriteFiles() []generator.WriteFile {
	return []generator.WriteFile{
		{
			FileName:     "types_on_off.hpp",
			TemplateName: "types_on_off.hpp",
		},
		{
			FileName:     "types_on_off.cpp",
			TemplateName: "types_on_off.cpp",
		},
	}
}

func (*OnOff) Includes() []string {
	return []string{
		"zephyr/drivers/gpio.h",
		"types_on_off.hpp",
	}
}

func (*OnOff) AppConfig() []appconfig.ConfigValue {
	return []appconfig.ConfigValue{
		appconfig.NewValue("CONFIG_GPIO").Required(appconfig.Yes),
	}
}

func (o *OnOff) ApplyOverlay(overlay *devicetree.DeviceTree) error {
	dtPin := devicetree.NewLED(o.Pin)
	return dtPin.AttachSelf(overlay)
}

func (*OnOff) Extenders() []generator.Extender {
	return []generator.Extender{
		extenders.GPIO{},
		extenders.NewLEDs(),
	}
}
