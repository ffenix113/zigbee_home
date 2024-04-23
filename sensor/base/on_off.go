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
}

func (*OnOff) String() string {
	return "On/Off"
}

func (*OnOff) Template() string {
	return "sensors/on_off"
}

func (o *OnOff) Clusters() cluster.Clusters {
	return []cluster.Cluster{
		cluster.OnOff{PinLabel: o.Pin.Label()},
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
	}
}
