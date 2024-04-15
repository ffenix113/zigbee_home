package sensor_test

import (
	"testing"

	"github.com/ffenix113/zigbee_home/sensor/bosch"
	"github.com/ffenix113/zigbee_home/types/sensor"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestSensorsUnmarshal(t *testing.T) {
	bts := []byte(`
- type: bme280
  connection:
    type: 'i2c1'
    address: '0x76'`)

	var sensors sensor.Sensors

	require.NoError(t, yaml.Unmarshal(bts, &sensors))

	require.Equal(t, "bme280", sensors[0].(*bosch.BME280).Base.Type)
	require.Equal(t, "0x76", sensors[0].(*bosch.BME280).Connection["address"])
}
