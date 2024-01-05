package base

import (
	"github.com/ffenix113/zigbee_home/cli/templates/extenders"
	"github.com/ffenix113/zigbee_home/cli/types"
	"github.com/ffenix113/zigbee_home/cli/types/appconfig"
	"github.com/ffenix113/zigbee_home/cli/types/devicetree"
	"github.com/ffenix113/zigbee_home/cli/types/generator"
	"github.com/ffenix113/zigbee_home/cli/zcl/cluster"
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
	dtPin := devicetree.Pin{Pin: o.Pin}
	return dtPin.AttachSelf(overlay)
}

func (*OnOff) Extenders() []generator.Extender {
	return []generator.Extender{
		extenders.GPIO{},
	}
}
