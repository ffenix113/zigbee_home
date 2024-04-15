Supporting new sensors can be really easy, or somewhat hard, depending on if the sensor has drivers in Zephyr and if all necessary Zigbee clusters are defined.

# Sensor base
All sensors must implement an interface `sensor.Sensor` defined in `types/sensor` package:
```go
type Sensor interface {
	// Stringer is for human-readable name
	fmt.Stringer
	// Label returns unique label value for the sensor.
	// Generally this method should not be defined by user,
	// intstead it will be defined in embedded `*base.Base`
	Label() string
	Template() string
	cluster.Provider
	appconfig.Provider
	devicetree.Applier
}
```

Each sensor must embed `*base.Base`, which is defined in `/sensor/base`, which provides all of the methods above as no-ops:
```go
type SomeSensor struct {
    *base.Base `yaml:",inline"` // The tag is required as well
}
```


Sensor can define

* Zigbee clusters - necessary for any sensor that should output values to coordinator, be controlled by coordinator, or both.
* Devicetree overlay configuration - should be necessary by all sensors, as devicetree defines how sensor is attached to the board.
* Configuration options - generally also would be required, as Zephyr may choose to include driver for the sensor based on this configuration.
* Template - if sensor is supported by Zephyr and all Zigbee clusters are defined - this should be optional, as known sensors will re-use common functionality.
* Extenders - if sensor requires additional functionality like adding required headers, or adding a driver via Zephyr module - this can be defined here.

# User-configurable and internal sensor options

Some sensors might need some user configuration to be working properly.

For example internal device temperature sensor does not need user configuration, as it will be defined the same for any board and peripherals configuration.

Others might need some input to input address of I2C sensor, set operating mode or disabling part of the sensor.

In some cases configuration options may also be used for internal purposes, like diferentiating between variants of the sensor. For this user should not be able to provide this configuration manually or change it in any way.

Yaml tag is used to tell if option is user-configurable or not. `yaml:"-"` would mean internal option, while omitting the tag, or specifying a field name like `yaml:"sample_rate"` would make it user-configurable.

This can be done by adding necessary configuration fields to sensor struct:
```go
type SomeSensor struct {
    *base.Base `yaml:",inline"`
    // All of fields below will be a part of sensor configuration.
    
    // This fields are internal, 
    // meaning that they will not be user-configurable.
    // Tag `yaml:"-"` means that this options will 
    // not be picked up from configuration file.
    Variant string `yaml:"-"`
    Frequency uint `yaml:"-"`
    
    // This options are user-configurable.
    // In configuration file they can be added by
    // lowercasing field name, or providing a name
    // as defined in `yaml` tag(if any).
    // Note that any Yaml unmarshalable type is supported.
    I2C base.I2CConnection // This will be `i2c` in configuration file
    OperationMode string // This will be `operationmode`
    SampleRate uint `yaml:"sample_rate"` // This will be `sample_rate`
}
```

# Adding sensor definition

All sensors have to be defined and added to a known sensor list.

Sensor definition provides a necessary information on how to add the sensor to the firmware. For example configuration options that should be added.

Sensor definition should be added to
```
/sensor/<manufacturer>/<sensor_name>.go
```
And a known sensor entry should be added to
```
/types/sensor/known.go
```

# Adding new Zigbee cluster
TBD

# Sensor supported by Zephyr
This is the easiest case, as only sensor definition is required, without any template files or custom drivers.

## Example: Bosch BME280
```go
package bosch

import (
	"strings"

	"github.com/ffenix113/zigbee_home/sensor/base"
	"github.com/ffenix113/zigbee_home/templates/extenders"
	"github.com/ffenix113/zigbee_home/types/appconfig"
	dt "github.com/ffenix113/zigbee_home/types/devicetree"
	"github.com/ffenix113/zigbee_home/types/generator"
	"github.com/ffenix113/zigbee_home/zcl/cluster"
)

type BME280 struct {
    // Required embedding of `*base.Base`
	*base.Base `yaml:",inline"`
    // This sensor works through I2C, 
    // so get the configuration for it
	I2C        base.I2CConnection
    // Variant is "internal" option, 
    // that defines if 280 or 680 is actually used
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
    // Cluster values below are taken from sensor datasheet.
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
		cluster.RelativeHumidity{
			MinMeasuredValue: 10,
			MaxMeasuredValue: 90,
		},
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
			dt.NewProperty("reg", dt.Angled(b.I2C.Reg())),
			dt.PropertyStatusEnable,
		},
	})

	return nil
}

// For now this Extenders is required for sensors,
// as it includes necessary configuration for Zephyr 
// to add sensor drivers and functions.
func (BME280) Extenders() []generator.Extender {
	return []generator.Extender{
		extenders.NewSensor(),
	}
}

```

# Sensor that are not supported by Zephyr
## Sensor without need for driver
TBD
## External driver as module
TBD