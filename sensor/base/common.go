package base

import (
	"github.com/ffenix113/zigbee_home/zcl/cluster"
)

type CommonDeviceClusters struct {
	*Base `yaml:",inline"`
}

func NewCommonDeviceClusters() *CommonDeviceClusters {
	return &CommonDeviceClusters{
		Base: &Base{
			label: "common_device_clusters",
		},
	}
}

func (*CommonDeviceClusters) String() string {
	return "Common device clusters"
}

func (*CommonDeviceClusters) Template() string {
	return "sensors/common"
}

func (o *CommonDeviceClusters) Clusters() cluster.Clusters {
	return []cluster.Cluster{
		cluster.Basic{},
	}
}
