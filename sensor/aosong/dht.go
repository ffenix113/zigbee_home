package aosong

import (
	"fmt"

	"github.com/ffenix113/zigbee_home/sensor/base"
	"github.com/ffenix113/zigbee_home/templates/extenders"
	"github.com/ffenix113/zigbee_home/types/appconfig"
	dt "github.com/ffenix113/zigbee_home/types/devicetree"
	"github.com/ffenix113/zigbee_home/types/generator"
	"github.com/ffenix113/zigbee_home/zcl/cluster"
	"github.com/ffenix113/zigbee_home/types"
)

type DHT struct {
	*base.Base `yaml:",inline"`
	Pin      types.Pin
    Variant  string `yaml:"variant"`
}

func NewDHT() *DHT {
	return &DHT{
		Variant: "",
	}
}

func (b DHT) String() string {
	return "Aosong"
}

func (DHT) Clusters() cluster.Clusters {
	return []cluster.Cluster{
		cluster.Temperature{
			MinMeasuredValue: -40,
			MaxMeasuredValue: 80,
			Tolerance:        1,
		},
        cluster.NewRelativeHumidity(0, 100),
	}
}

func (b DHT) AppConfig() []appconfig.ConfigValue {
	return []appconfig.ConfigValue{
		appconfig.NewValue("CONFIG_GPIO").Required(appconfig.Yes),
		appconfig.CONFIG_DHT.Required(appconfig.Yes),
	}
}

func (b DHT) ApplyOverlay(tree *dt.DeviceTree) error {

    pinctrlNode := tree.FindSpecificNode(dt.SearchByLabel(dt.NodeLabelPinctrl))
	if pinctrlNode == nil {
		return dt.ErrNodeNotFound(dt.NodeLabelPinctrl)
	}

    pinLabel := fmt.Sprintf("gpio%d %d (GPIO_PULL_UP | GPIO_ACTIVE_LOW)", b.Pin.Port, b.Pin.Pin)

    props:= []dt.Property{
        dt.NewProperty("compatible", dt.FromValue("aosong,dht")),
        dt.NewProperty("status", dt.FromValue("okay")),
        dt.NewProperty("dio-gpios", dt.Angled(dt.Label(pinLabel))),
    }
    // Add variant property only if Variant is not emtpy
    // For dht22 or AM2302 devices the string "dht22;" must be added to device tree overlay
    // See : https://docs.zephyrproject.org/latest/build/dts/api/bindings/sensor/aosong,dht.html
    if b.Variant != "" {
        props = append(props, dt.NewProperty(b.Variant, nil))
    }

    pinctrlNode.AddNodes(
        &dt.Node{
            Name:  "dht22",
            Label: b.Label(),
            Properties: props,
        },
    )

	return nil
}

func (DHT) Extenders() []generator.Extender {
	return []generator.Extender{
		extenders.NewSensor(),
		extenders.GPIO{},
	}
}
