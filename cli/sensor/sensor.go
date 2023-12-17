package sensor

import (
	"fmt"
	"reflect"

	"github.com/ffenix113/zigbee_home/cli/generate/appconfig"
	"github.com/ffenix113/zigbee_home/cli/generate/devicetree"
	"github.com/ffenix113/zigbee_home/cli/sensor/base"
	"github.com/ffenix113/zigbee_home/cli/zcl/cluster"
	"gopkg.in/yaml.v3"
)

type Sensors []Sensor

type Sensor interface {
	Clusters() cluster.Clusters
	AppConfig() []appconfig.ConfigValue
	ApplyOverlay(overlay *devicetree.DeviceTree) error
}

func (s *Sensors) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.SequenceNode {
		return fmt.Errorf("must have sequence, but have %q", value.Kind)
	}

	for i, node := range value.Content {
		sensor, err := unmarshalSensor(node)
		if err != nil {
			return fmt.Errorf("unmarshal sensor %d: %w", i, err)
		}

		(*s) = append((*s), sensor)
	}

	return nil
}

func unmarshalSensor(node *yaml.Node) (Sensor, error) {
	var sensorType base.SensorType
	if err := node.Decode(&sensorType); err != nil {
		return nil, fmt.Errorf("get sensor type: %w", err)
	}

	sensorConfig, ok := knownSensors[sensorType.Type]
	if !ok {
		return nil, fmt.Errorf("unsupported sensor type: %q", sensorType.Type)
	}

	rVal := reflect.New(reflect.TypeOf(sensorConfig).Elem())
	if err := node.Decode(rVal.Interface()); err != nil {
		return nil, fmt.Errorf("decode sensor type %q: %w", sensorType.Type, err)
	}

	return rVal.Elem().Interface().(Sensor), nil
}
