package base

import (
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
	Button   string              `yaml:"button"`
	ZoneType cluster.IasZoneType `yaml:"zone_type"`
}

func (*IASZone) String() string {
	return "IAS Zone"
}

func (*IASZone) Template() string {
	return "sensors/ias_zone"
}

func (*IASZone) CPPComponentType() string {
	return "IASZone"
}

func (*IASZone) WriteFiles() []generator.WriteFile {
	return []generator.WriteFile{
		{
			FileName:     "types_ias_zone.hpp",
			TemplateName: "types_ias_zone.hpp",
		},
		{
			FileName:     "types_ias_zone.cpp",
			TemplateName: "types_ias_zone.cpp",
		},
	}
}

func (z *IASZone) Clusters() cluster.Clusters {
	// By default - be contact sensor for now.
	if z.ZoneType == "" {
		z.ZoneType = cluster.IasZoneContact
	}

	return []cluster.Cluster{
		cluster.IASZone{ZoneType: z.ZoneType},
	}
}
