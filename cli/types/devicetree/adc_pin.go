package devicetree

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/ffenix113/zigbee_home/cli/types"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

var resolutions = []uint8{8, 10, 12, 14}
var pinConversionMap = map[uint8]string{
	2:  "AIN0",
	3:  "AIN1",
	4:  "AIN2",
	5:  "AIN3",
	28: "AIN4",
	29: "AIN5",
	30: "AIN6",
	31: "AIN7",
}

var referenecs = map[string]string{
	"vdd_1_4":  "VDD_1_4",
	"internal": "INTERNAL",
}

var aquisitionTimeUnits = map[string]string{
	"ms":    "ADC_ACQ_TIME_MICROSECONDS",
	"ns":    "ADC_ACQ_TIME_NANOSECONDS",
	"ticks": "ADC_ACQ_TIME_TICKS",
}

type aquisitionTime struct {
	Value uint16
	Unit  string
}

func (at aquisitionTime) String() string {
	if at.Value == 0 {
		return "ADC_ACQ_TIME_DEFAULT"
	}

	return fmt.Sprintf("ADC_ACQ_TIME(%s, %d)", aquisitionTimeUnits[at.Unit], at.Value)
}

type ADCPin struct {
	// Define configurations,
	// as not all usages will be equal
	Gain           string
	Reference      string
	Resolution     uint8
	Oversampling   uint8
	AquisitionTime aquisitionTime

	Pin types.Pin
}

func (p ADCPin) Name() string {
	return p.Pin.Name()
}

func (p *ADCPin) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		return p.Pin.UnmarshalYAML(node)
	}

	type a ADCPin
	return node.Decode((*a)(p))
}

func (p ADCPin) AttachSelf(dt *DeviceTree) error {
	pinName := p.Pin.Name()

	p.setDefaults()
	if err := p.validate(); err != nil {
		return fmt.Errorf("validate adc pin: %w", err)
	}

	const nodeNameADC = "adc"
	adcNode := dt.FindSpecificNode(SearchByName(nodeNameADC))
	if adcNode == nil {
		adcNode = &Node{
			Label:  nodeNameADC,
			Upsert: true,
			Properties: []Property{
				NewProperty("#address-cells", FromValue(1)),
				NewProperty("#size-cells", FromValue(0)),
			},
		}

		dt.AddNodes(adcNode)
	}
	// We can safely do this because we check
	// if the pin is correct one on validation step.
	positivePinName, _ := pinConversionMap[p.Pin.Pin.Value()]

	numericLabel := p.Pin.NumericLabel()
	adcNode.AddNodes(&Node{
		Name:        pinName,
		UnitAddress: numericLabel,
		Properties: []Property{
			NewProperty("reg", Angled(String(numericLabel))),
			NewProperty("zephyr,gain", Quoted("ADC_GAIN_"+p.Gain)),
			NewProperty("zephyr,reference", Quoted("ADC_REF_"+referenecs[p.Reference])),
			NewProperty("zephyr,acquisition-time", Angled(String(p.AquisitionTime.String()))),
			NewProperty("zephyr,input-positive", Angled(String("NRF_SAADC_"+positivePinName))),
			NewProperty("zephyr,oversampling", FromValue(p.Oversampling)),
			NewProperty("zephyr,resolution", FromValue(p.Resolution)),
		},
	})

	zephyrUserNode := "zephyr,user"
	zephyrUser := dt.FindSpecificNode(SearchByName(NodeNameRoot), SearchByName(zephyrUserNode))
	if zephyrUser == nil {
		zephyrUser = &Node{
			Name: zephyrUserNode,
		}

		dt.FindSpecificNode(SearchByName(NodeNameRoot)).AddNodes(zephyrUser)
	}

	zephyrUser.Properties = append(zephyrUser.Properties,
		NewProperty("io-channels", Angled(String("&adc "+numericLabel))),
		NewProperty("io-channel-names", FromValue(pinName)),
	)

	return nil
}

func (p *ADCPin) setDefaults() {
	if p.Gain == "" {
		p.Gain = "1_6"
	}

	if p.Reference == "" {
		p.Reference = "internal"
	}

	if p.Resolution == 0 {
		p.Resolution = 12
	}
}

func (p ADCPin) validate() error {
	// If pin is not defined - do not add its configuration.
	if !p.Pin.PinsDefined() {
		return fmt.Errorf("must define pins to use, pin label: %q", p.Pin.Name())
	}

	if !p.Pin.Port.HasValue() || p.Pin.Port.Value() != 0 {
		return errors.New("port must always be 0")
	}

	if _, ok := pinConversionMap[p.Pin.Pin.Value()]; !ok {
		return fmt.Errorf("pin %d cannot be used as ADC pin", p.Pin.Pin.Value())
	}

	if _, ok := referenecs[p.Reference]; !ok {
		return fmt.Errorf("reference value is invalid: %q, valid values are: %v", p.Reference, strings.Join(maps.Keys(referenecs), ", "))
	}

	if !slices.Contains(resolutions, p.Resolution) {
		return fmt.Errorf("resolution is invalid: %d, valid values are: %v", p.Resolution, resolutions)
	}

	// A bit "hacky" way to check, because we want to fit
	if p.Oversampling > 8 {
		return fmt.Errorf("oversampling value is larger than max: %d > %d", p.Oversampling, 8)
	}

	if err := p.validateAquisitionTime(); err != nil {
		return fmt.Errorf("aquisition time is invalid: %w", err)
	}

	return nil
}

func (p ADCPin) validateAquisitionTime() error {
	aqTime := p.AquisitionTime

	if aqTime.Value == 0 {
		// This will be default, nothing to do here.
		return nil
	}

	var validUnit bool
	for unit := range aquisitionTimeUnits {
		if strings.HasSuffix(aqTime.Unit, unit) {
			validUnit = true
			break
		}
	}
	if !validUnit {
		return errors.New("invalid unit")
	}

	if maxAqTime := uint16(8191); aqTime.Value >= maxAqTime {
		return fmt.Errorf("aquisition time is too large: %d > %d", aqTime.Value, maxAqTime)
	}

	return nil
}
