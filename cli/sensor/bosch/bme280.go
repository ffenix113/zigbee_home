package bosch

import (
	"strings"

	"github.com/ffenix113/zigbee_home/cli/sensor/base"
	"github.com/ffenix113/zigbee_home/cli/templates/extenders"
	"github.com/ffenix113/zigbee_home/cli/types/appconfig"
	dt "github.com/ffenix113/zigbee_home/cli/types/devicetree"
	"github.com/ffenix113/zigbee_home/cli/types/generator"
	"github.com/ffenix113/zigbee_home/cli/zcl/cluster"
)

type BME280 struct {
	*base.Base `yaml:",inline"`
	I2C        base.I2CConnection
	Variant    string `yaml:"-"`
}

func NewBME280() *BME280 {
	return &BME280{
		Variant: "bme280",
	}
}

func (b BME280) String() string {
	return "Bosch " + strings.ToUpper(b.Variant)
}

func (BME280) Clusters() cluster.Clusters {
	return []cluster.Cluster{
		cluster.Temperature{
			MinMeasuredValue: -40,
			MaxMeasuredValue: 85,
			Tolerance:        1,
		},
		cluster.Pressure{
			MinMeasuredValue: 30,
			MaxMeasuredValue: 110,
			Tolerance:        0,
		},
		cluster.NewRelativeHumidity(10, 90),
	}
}

func (b BME280) AppConfig() []appconfig.ConfigValue {
	return []appconfig.ConfigValue{
		appconfig.CONFIG_I2C.Required(appconfig.Yes),
		// Yes, for both 280 & 680 we are setting through BME280
		appconfig.CONFIG_BME280.Required(appconfig.Yes),
		appconfig.NewValue("CONFIG_BME280_MODE_FORCED").Required(appconfig.Yes),
	}
}

func (b BME280) ApplyOverlay(tree *dt.DeviceTree) error {
	i2cNode := tree.FindSpecificNode(dt.SearchByLabel(b.I2C.ID))
	if i2cNode == nil {
		return dt.ErrNodeNotFound(b.I2C.ID)
	}

	i2cNode.AddNodes(&dt.Node{
		Name:        b.Variant,
		Label:       b.Label(),
		UnitAddress: b.I2C.UnitAddress(),
		Properties: []dt.Property{
			dt.NewProperty(dt.PropertyNameCompatible, dt.FromValue("bosch,"+b.Variant)),
			dt.NewProperty("reg", dt.Angled(dt.String(b.I2C.Reg()))),
			dt.PropertyStatusEnable,
		},
	})

	return nil
}

func (BME280) Extenders() []generator.Extender {
	return []generator.Extender{
		extenders.NewSensor(),
	}
}
