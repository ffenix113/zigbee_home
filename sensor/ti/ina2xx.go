package ti

import (
	"fmt"
	"math"
	"math/bits"
	"strconv"
	"strings"

	"github.com/ffenix113/zigbee_home/sensor/base"
	"github.com/ffenix113/zigbee_home/templates/extenders"
	dt "github.com/ffenix113/zigbee_home/types/devicetree"
	"github.com/ffenix113/zigbee_home/types/generator"
	"github.com/ffenix113/zigbee_home/zcl/cluster"
)

type INA2XX struct {
	*base.Base `yaml:",inline"`

	Model string

	I2C base.I2CConnection

	// Max current to be measured in Amps.
	MaxCurrentAmps float32 `yaml:"max_current_amps"`
	// Shunt resistance in milli Ohms.
	ShuntResistanceMohms float32 `yaml:"shunt_resistance_mohms"`
	// Number of samples to average.
	Avg uint16
}

func (i *INA2XX) String() string {
	return fmt.Sprintf("INA2XX (%q)", strings.ToLower(i.Model))
}

func (*INA2XX) CPPComponentType() string {
	return "BasicSensor"
}

func (*INA2XX) Clusters() cluster.Clusters {
	return []cluster.Cluster{
		cluster.ElectricalMeasurement{},
	}
}

func (*INA2XX) NeedsDevice() bool {
	return true
}

func (i *INA2XX) ApplyOverlay(overlay *dt.DeviceTree) error {
	if i.Model == "" {
		return fmt.Errorf("model cannot be empty")
	}

	i.Model = strings.ToLower(i.Model)

	if i.I2C.ID == "" {
		return fmt.Errorf("i2c id cannot be empty")
	}

	if i.ShuntResistanceMohms <= 0 {
		return fmt.Errorf("shunt resister value cannot be less or equal to 0")
	}

	if i.MaxCurrentAmps <= 0 {
		return fmt.Errorf("max current cannot be less or equal to 0")
	}

	if i.Avg == 0 {
		i.Avg = 1
	}

	if i.Avg > 1024 {
		i.Avg = 1024
	}

	if bits.OnesCount16(i.Avg) != 1 {
		return fmt.Errorf("avg count should be power of two, but was %d", i.Avg)
	}

	aliases := overlay.FindSpecificNode(dt.SearchByName(dt.NodeNameRoot), dt.SearchByName(dt.NodeNameAliases))
	aliases.AddProperties(dt.NewProperty(i.DeviceTreeLabel(), dt.Label(i.Label())))

	i2c := overlay.FindSpecificNode(dt.SearchByLabel(i.I2C.ID))
	if i2c == nil {
		return fmt.Errorf("i2c bus with id %q should be already set up", i.I2C.ID)
	}

	if i.I2C.Addr == 0 {
		// Default for ina226.
		// Yes, this may be wrong for other models or configurations.
		i.I2C.Addr = 0x40
	}

	// TODO: verify this value
	shuntMicroOhmsStr := strconv.Itoa(int(i.ShuntResistanceMohms * 1_000))
	// TODO: verify this value
	currentMicroAmpsLsb := int(i.MaxCurrentAmps * 1_000_000 / (math.MaxInt16 + 1))

	// 1 lsb step in microvolts.
	LsbUvForModel := map[string]float32{
		"ina219": 10,     // 10 uV
		"ina226": 2.5,    // 2.5 uV
		"ina228": 0.3125, // 312.5 nV, ADCRANGE = 0
		"ina230": 2.5,    // 2.5 uV
		"ina232": 2.5,    // 2.5 uV, ADCRANGE = 0
		"ina236": 2.5,    // 2.5 uV, ADCRANGE = 0
		"ina237": 5,      // 5 uV, ADCRANGE = 0
	}

	uVForModel := LsbUvForModel[i.Model]
	if uVForModel == 0 {
		return fmt.Errorf("unsupported ina2xx model to fetch min uV per LSB: %q", i.Model)
	}

	// Verify that we did not fall through the floor of possible value of the chip.
	// lsbUv / Rshunt > currentMicroAmpsLSB
	// 10^-6 / 10^-3 == 10^-3
	//
	// This calculation is done explicitely, without simplifying
	// so it is a bit more explicit.
	minAmps := (uVForModel / 1_000_000) / (i.ShuntResistanceMohms / 1_000)
	if minMicroAmpsLsb := 1_000_000 * minAmps; minMicroAmpsLsb > float32(currentMicroAmpsLsb) {
		currentMicroAmpsLsb = int(minMicroAmpsLsb)
	}

	currentMicroAmpsStr := strconv.Itoa(currentMicroAmpsLsb)

	inaNode := &dt.Node{
		Name:        i.Model,
		Label:       i.Label(),
		UnitAddress: i.I2C.UnitAddress(),
		Properties: []dt.Property{
			dt.PropertyStatusEnable,
			dt.NewProperty(dt.PropertyNameCompatible, dt.Quoted("ti,"+i.Model)),

			dt.NewProperty("reg", dt.Angled(dt.String(i.I2C.Reg()))),
		},
	}

	var modelProps []dt.Property
	// INA219 needs some other (older?) properties
	if i.Model == "ina219" {
		modelProps = i.dtProperties_ina219(shuntMicroOhmsStr, currentMicroAmpsStr)
	} else {
		modelProps = i.dtProperties_ina2xx(shuntMicroOhmsStr, currentMicroAmpsStr)
	}

	i2c.AddNodes(inaNode.AddProperties(modelProps...))

	return nil
}

func (*INA2XX) Extenders() []generator.Extender {
	return []generator.Extender{
		extenders.NewSensor(),
		extenders.ElectricalMeasurement{},
	}
}

func (i *INA2XX) dtProperties_ina219(shuntUOhms, currentUAmps string) []dt.Property {
	return []dt.Property{
		dt.NewProperty("shunt-milliohm", dt.Angled(dt.String(shuntUOhms))),
		dt.NewProperty("lsb-microamp", dt.Angled(dt.String(currentUAmps))),
	}
}

func (i *INA2XX) dtProperties_ina2xx(shuntUOhms, currentUAmps string) []dt.Property {
	return []dt.Property{
		dt.NewProperty("avg-count", dt.FromValue(i.Avg)),
		dt.NewProperty("rshunt-micro-ohms", dt.Angled(dt.String(shuntUOhms))),
		dt.NewProperty("current-lsb-microamps", dt.Angled(dt.String(currentUAmps))),
	}
}
