package types

import (
	"github.com/ffenix113/zigbee_home/cli/config"
	"github.com/ffenix113/zigbee_home/cli/zcl/cluster"
)

type TemplateDevice struct {
	Config *config.Device

	Endpoints Endpoints
}

func (d *TemplateDevice) AddCluster(cluster cluster.Cluster) error {
	endpoint := d.Endpoints.LastEndpoint()

	if endpoint.HasCluster(cluster.ID()) {
		endpoint = d.Endpoints.AddEndpoint()
	}

	endpoint.Clusters = append(endpoint.Clusters, cluster)

	return nil
}
