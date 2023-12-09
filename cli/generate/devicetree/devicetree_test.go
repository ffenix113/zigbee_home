package devicetree_test

import (
	"strings"
	"testing"

	"github.com/ffenix113/zigbee_home/cli/generate/devicetree"
)

func TestWriteOverlay(t *testing.T) {
	deviceTree := (&devicetree.DeviceTree{}).AddNodes(
		devicetree.Node{
			Name: devicetree.NodeNameRoot,
			SubNodes: []devicetree.Node{
				{
					Name: devicetree.NodeNameChosen,
					Properties: []devicetree.Property{
						devicetree.NewProperty("ncs,zigbee-timer", devicetree.Label("timer2")),
						devicetree.NewProperty("zephyr,console", devicetree.Label("cdc_acm_uart0")),
						devicetree.NewProperty("zephyr,entropy", devicetree.Label("rng")),
					},
				},
			},
		},
		devicetree.Node{
			Name: "&zephyr_udc0",
			SubNodes: []devicetree.Node{
				{
					Name:  "cdc_acm_uart0",
					Label: "cdc_acm_uart0",
					Properties: []devicetree.Property{
						devicetree.NewProperty(devicetree.PropertyNameCompatible, devicetree.String("zephyr,cdc-acm-uart")),
					},
				},
			},
		},
	)

	var buf strings.Builder

	deviceTree.WriteTo(&buf)

	t.Log(buf.String())
}
