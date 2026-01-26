package sensor

import (
	"fmt"
	"reflect"

	"github.com/ffenix113/zigbee_home/sensor"
	"github.com/ffenix113/zigbee_home/sensor/aosong"
	"github.com/ffenix113/zigbee_home/sensor/base"
	"github.com/ffenix113/zigbee_home/sensor/bosch"
	"github.com/ffenix113/zigbee_home/sensor/sensirion"
	"github.com/ffenix113/zigbee_home/sensor/ti"
)

func fromType[T Sensor]() Sensor {
	var s T

	rVal := reflect.New(reflect.TypeOf(s).Elem())

	return rVal.Interface().(Sensor)
}

func fromConstructor(constr any) func() Sensor {
	return func() Sensor {
		rVal := reflect.ValueOf(constr)

		if numOut := rVal.Type().NumOut(); numOut != 1 {
			panic(fmt.Sprintf("constructor should return exactly 1 value, but has %d", numOut))
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

	"rotary_encoder": fromConstructor(base.NewRotaryEncoder),

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

	// Aosong
	"dht": fromConstructor(aosong.NewDHT),

	// Sensirion
	"scd4x": fromType[*sensirion.SCD4X],

	// Texas Instruments
	"ina2xx": fromType[*ti.INA2XX],
}
