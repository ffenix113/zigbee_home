package sensirion

import (
	"github.com/ffenix113/zigbee_home/cli/sensor/base"
	"github.com/ffenix113/zigbee_home/cli/templates/extenders"
	"github.com/ffenix113/zigbee_home/cli/types/appconfig"
	dt "github.com/ffenix113/zigbee_home/cli/types/devicetree"
	"github.com/ffenix113/zigbee_home/cli/types/generator"
	"github.com/ffenix113/zigbee_home/cli/zcl/cluster"
)

type SCD4X struct {
	*base.Base        `yaml:",inline"`
	I2C               base.I2CConnection
	TemperatureOffset int8 `yaml:"temperature_offset"`
}

func (SCD4X) String() string {
	// TODO: Update to SCD4X when will support properly SCD40
	return "Sensirion SCD4X (SCD41)"
}

func (SCD4X) Clusters() cluster.Clusters {
	// https://sensirion.com/media/documents/E0F04247/631EF271/CD_DS_SCD40_SCD41_Datasheet_D1.pdf
	return []cluster.Cluster{
		cluster.Temperature{
			MinMeasuredValue: -10,
			MaxMeasuredValue: 60,
			Tolerance:        1,
		},
		cluster.NewRelativeHumidity(0, 100),
		cluster.CarbonDioxide{
			MinMeasuredValue: 400,
			MaxMeasuredValue: 5000,
			Tolerance:        40,
		},
	}
}

func (SCD4X) AppConfig() []appconfig.ConfigValue {
	return []appconfig.ConfigValue{
		appconfig.CONFIG_I2C.Required(appconfig.Yes),
		appconfig.NewValue("CONFIG_CRC").Required(appconfig.Yes),
		appconfig.NewValue("CONFIG_SCD4X").Required(appconfig.Yes),
	}
}

func (s SCD4X) ApplyOverlay(tree *dt.DeviceTree) error {
	i2cNode := tree.FindSpecificNode(dt.SearchByLabel(s.I2C.ID))
	if i2cNode == nil {
		return dt.ErrNodeNotFound(s.I2C.ID)
	}
	// SCD4X address is static
	s.I2C.Addr = "0x62"

	i2cNode.AddNodes(&dt.Node{
		Name:        "scd4x",
		Label:       s.Label(),
		UnitAddress: s.I2C.UnitAddress(),
		Properties: []dt.Property{
			dt.PropertyStatusEnable,
			dt.NewProperty(dt.PropertyNameCompatible, dt.FromValue("sensirion,scd4x")),
			dt.NewProperty("reg", dt.Angled(dt.String(s.I2C.Reg()))),
			// Only single-shot for now.
			// Would need some changes in templates for changing
			dt.NewProperty("measure-mode", dt.FromValue("single-shot")),
			dt.NewProperty("model", dt.FromValue("scd41")),
			dt.NewProperty("temperature-offset", dt.FromValue(s.TemperatureOffset)),
		},
	})

	return nil
}

func (s SCD4X) Extenders() []generator.Extender {
	return []generator.Extender{
		extenders.NewSensor(),
		extenders.NewSCD4X(),
	}
}
