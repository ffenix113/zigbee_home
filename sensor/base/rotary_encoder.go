package base

import (
	"strconv"
	"time"

	"github.com/ffenix113/zigbee_home/templates/extenders"
	"github.com/ffenix113/zigbee_home/types"
	"github.com/ffenix113/zigbee_home/types/appconfig"
	dt "github.com/ffenix113/zigbee_home/types/devicetree"
	"github.com/ffenix113/zigbee_home/types/generator"
	"github.com/ffenix113/zigbee_home/zcl/cluster"
)

// Rotary encoder will send events for each rotation so that devices,
// such as lights, can be level-controlled, for example the brightness.
//
// Rotary encoder can also optionally have a switch (button).
// Though currently switch is no supported.
//
// TODO: Add an option to count total number of pulses(rotations)
// instead of emitting after each (debounced) pulse.
//
// So instead of sending a "step" command every half a second,
// it could count all pulses, and send a cumulative "step" command
// after some period of time, for example 1 second.
type RotaryEncoder struct {
	*Base                 `yaml:",inline"`
	MaxRotationsPerSecond uint8 `yaml:"max_rotations_per_second"`
	// Resolution optionally defines how many pulses should be registered before reporting the change.
	//
	// If rotary encoder is just left/right -
	// then most probably this value should be left with default value (2).
	Resolution uint8 `yaml:"resolution"`
	// Sample time optionally specifies how frequently pins should be polled when first pulse is received.
	// The lower this value - higher the rotation speed can be used.
	//
	// Probably it is not reasonable to set this value to more than 2ms.
	//
	// If you notice that a lot of steps are not reported - try lowering this value by 100ms
	// until you are satisfied with reported counts.
	SampleTime time.Duration `yaml:"sample_time"`

	Pins struct {
		A types.PinWithID
		B types.PinWithID

		SwitchPin types.PinWithID `yaml:"switch_pin"`
	}
}

func NewRotaryEncoder() *RotaryEncoder {
	return &RotaryEncoder{
		MaxRotationsPerSecond: 2,
		Resolution:            2,
		SampleTime:            1500 * time.Microsecond,
	}
}

func (*RotaryEncoder) String() string {
	return "RotaryEncoder"
}

func (*RotaryEncoder) Template() string {
	return "sensors/level_control"
}

func (*RotaryEncoder) CPPComponentType() string {
	return "LevelControl"
}

func (o *RotaryEncoder) Clusters() cluster.Clusters {
	return []cluster.Cluster{
		cluster.LevelControl{},
	}
}

func (*RotaryEncoder) WriteFiles() []generator.WriteFile {
	return []generator.WriteFile{
		{
			FileName:     "types_level_control.hpp",
			TemplateName: "types_level_control.hpp",
		},
		{
			FileName:     "types_level_control.cpp",
			TemplateName: "types_level_control.cpp",
		},
	}
}

func (*RotaryEncoder) Includes() []string {
	return []string{
		"types_level_control.hpp",
	}
}

func (*RotaryEncoder) AppConfig() []appconfig.ConfigValue {
	return []appconfig.ConfigValue{
		appconfig.NewValue("CONFIG_GPIO").Required(appconfig.Yes),
		appconfig.NewValue("CONFIG_INPUT").Required(appconfig.Yes),
	}
}

func (e *RotaryEncoder) NeedsDevice() bool {
	return true
}

func (o *RotaryEncoder) ApplyOverlay(overlay *dt.DeviceTree) error {
	resolution := 2
	if o.Resolution != 0 {
		resolution = int(o.Resolution)
	}

	sampleTime := 1500 * time.Microsecond
	if o.SampleTime != 0 {
		sampleTime = o.SampleTime
	}

	root := overlay.FindSpecificNode(dt.SearchByName(dt.NodeNameRoot))

	root.AddNodes(&dt.Node{
		Name:  o.Label(),
		Label: o.Label(),

		Properties: []dt.Property{
			dt.NewProperty(dt.PropertyNameCompatible, dt.Quoted("gpio-qdec")),
			dt.PropertyStatusEnable,
			dt.NewProperty("gpios", dt.Array(
				dt.GPIOPin(o.Pins.A.ToPin()),
				dt.GPIOPin(o.Pins.B.ToPin()),
			)),
			dt.NewProperty("steps-per-period", dt.Angled(dt.String(strconv.Itoa(resolution)))),

			dt.NewProperty("zephyr,axis", dt.Angled(dt.String("INPUT_REL_WHEEL"))),
			dt.NewProperty("sample-time-us", dt.Angled(dt.String(strconv.Itoa(int(sampleTime.Microseconds()))))),
			dt.NewProperty("idle-timeout-ms", dt.Angled(dt.String("400"))),
		},
	})

	aliases := overlay.FindSpecificNode(dt.SearchByName(dt.NodeNameRoot), dt.SearchByName(dt.NodeNameAliases))

	aliases.AddProperties(dt.NewProperty(o.DeviceTreeLabel(), dt.Label(o.Label())))

	// Currently switch is not supported.

	return nil
}

func (*RotaryEncoder) Extenders() []generator.Extender {
	return []generator.Extender{
		extenders.GPIO{},
		// extenders.NewLEDs(),
	}
}
