package base

import (
	"github.com/ffenix113/zigbee_home/templates/extenders"
	"github.com/ffenix113/zigbee_home/types"
	"github.com/ffenix113/zigbee_home/types/appconfig"
	"github.com/ffenix113/zigbee_home/types/generator"
	"github.com/ffenix113/zigbee_home/zcl/cluster"
)

func NewContact() *IASZone {
	return &IASZone{
		ZoneType: cluster.IasZoneContact,
	}
}

type IASZone struct {
	*Base    `yaml:",inline"`
	Pin      types.Pin
	ZoneType cluster.IasZoneType `yaml:"zone_type"`
}

func (*IASZone) String() string {
	return "IAS Zone"
}

func (*IASZone) Template() string {
	return "sensors/ias_zone"
}

func (o *IASZone) Clusters() cluster.Clusters {
	// By default - be contact sensor for now.
	if o.ZoneType == "" {
		o.ZoneType = cluster.IasZoneContact
	}

	return []cluster.Cluster{
		cluster.IASZone{ZoneType: o.ZoneType},
	}
}

func (*IASZone) AppConfig() []appconfig.ConfigValue {
	return []appconfig.ConfigValue{
		appconfig.NewValue("CONFIG_GPIO").Required(appconfig.Yes),
	}
}

func (z *IASZone) Extenders() []generator.Extender {
	return []generator.Extender{
		extenders.GPIO{},
		extenders.NewButtons(z.Pin),
	}
}
