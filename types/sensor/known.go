package sensor

import (
	"fmt"
	"reflect"

	"github.com/ffenix113/zigbee_home/sensor"
	"github.com/ffenix113/zigbee_home/sensor/base"
	"github.com/ffenix113/zigbee_home/sensor/bosch"
	"github.com/ffenix113/zigbee_home/sensor/sensirion"
)

func fromType[T Sensor]() Sensor {
	var s T

	rVal := reflect.New(reflect.TypeOf(s).Elem())

	return rVal.Interface().(Sensor)
}

func fromConstructor(constr any) func() Sensor {
	return func() Sensor {
		rVal := reflect.ValueOf(constr)

		numOut := rVal.Type().NumOut()
		switch {
		case numOut == 0:
			panic("constructor must have 1 return value")
		case numOut > 1:
			retType := rVal.Type().Out(0)
			panic(fmt.Sprintf("constructor %q should return exactly 1 value", retType.String()))
		}

		ret := rVal.Call(nil)[0]

		return ret.Interface().(Sensor)
	}
}

var knownSensors = map[string]func() Sensor{
	// Generic
	"on_off":       fromType[*base.OnOff],
	"power_config": fromType[*base.PowerConfiguration],
	"contact":      fromConstructor(base.NewContact),
	// Later we can just alias this to `soil_moisture`
	// if `soil_moisture` will not be used otherwise.
	"soil_moisture_adc": fromType[*base.SoilMoistureADC],
	// Generic ias zone sensor.
	// While it is defined here - for now it is
	// not useful much, as it only can be used
	// as contact sensor.
	"ias_zone": fromType[*base.IASZone],

	// Specific devices

	"device_temperature": fromType[*sensor.DeviceTemperature],

	// Bosch
	"bme280": fromConstructor(bosch.NewBME280),
	// This is a clone of bme280, with different overlay name
	// FIXME: It does not yet support IAQ measurements,
	// and does not expose resistance to Zigbee.
	"bme680": fromConstructor(bosch.NewBME680),

	// Sensirion
	"scd4x": fromType[*sensirion.SCD4X],
}
