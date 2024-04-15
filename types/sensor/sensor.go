package sensor

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/ffenix113/zigbee_home/sensor/base"
	"github.com/ffenix113/zigbee_home/types/appconfig"
	"github.com/ffenix113/zigbee_home/types/devicetree"
	"github.com/ffenix113/zigbee_home/types/generator"
	"github.com/ffenix113/zigbee_home/zcl/cluster"
	"gopkg.in/yaml.v3"
)

var sensorCounter int

// SensorLabelFn returns a unique label for a sensor.
var SensorLabelFn = func(s Sensor) string {
	sensorCounter += 1
	cleanLabel := strings.ReplaceAll(strings.ToLower(fmt.Sprintf("%T_", s)), ".", "_")
	return strings.TrimPrefix(cleanLabel, "*") + strconv.Itoa(sensorCounter)
}

type Sensors []Sensor

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

type WithExtenders interface {
	Extenders() []generator.Extender
}

// UniqueClusters will return all unique clusters across all the sensors.
// Clusters might be configured for a specific sensor,
// so this method is mostly useful to *know* which clusters are available.
func (s *Sensors) UniqueClusters() cluster.Clusters {
	clusterMap := map[cluster.ID]struct{}{}

	var clusters cluster.Clusters

	for _, sensor := range *s {
		for _, cluster := range sensor.Clusters() {
			_, ok := clusterMap[cluster.ID()]
			if ok {
				continue
			}

			clusterMap[cluster.ID()] = struct{}{}
			clusters = append(clusters, cluster)
		}
	}

	return clusters
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

	sensorConfigConstructor, ok := knownSensors[sensorType.Type]
	if !ok {
		return nil, fmt.Errorf("unsupported sensor type: %q", sensorType.Type)
	}

	rVal := reflect.ValueOf(sensorConfigConstructor())
	if err := node.Decode(rVal.Interface()); err != nil {
		return nil, fmt.Errorf("decode sensor type %q: %w", sensorType.Type, err)
	}

	sensor := rVal.Interface().(Sensor)
	base := rVal.Elem().FieldByName("Base").Interface().(*base.Base)
	base.SetLabel(SensorLabelFn(sensor))

	return sensor, nil
}
