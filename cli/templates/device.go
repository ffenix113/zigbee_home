package templates

import (
	"github.com/ffenix113/zigbee_home/cli/config"
	"github.com/ffenix113/zigbee_home/cli/types"
)

func ConfigToTemplateDevice(device *config.Device) *types.TemplateDevice {
	td := &types.TemplateDevice{
		RunEvery: device.General.RunEvery,
	}

	for _, sensor := range device.Sensors {
		sensorClusters := sensor.Clusters()

		endpoint := td.Endpoints.LastEndpoint()

		for _, cluster := range sensorClusters {
			if endpoint.HasCluster(cluster.ID()) {
				endpoint = td.Endpoints.AddEndpoint()
				break
			}
		}

		for _, cluster := range sensorClusters {
			endpoint.AddCluster(cluster)
		}
	}

	return td
}
