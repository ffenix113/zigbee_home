package extenders

import (
	"fmt"
	"path/filepath"

	"github.com/ffenix113/zigbee_home/types/devicetree"
	"github.com/ffenix113/zigbee_home/types/generator"
)

var _ generator.Extender = ADC{}
var _ devicetree.Applier = ADC{}

type ADC struct {
	generator.SimpleExtender

	Instances []devicetree.ADCPin
}

func NewADC(instances ...devicetree.ADCPin) generator.Extender {
	return ADC{
		Instances: instances,
	}
}

func (l ADC) Template() string {
	return filepath.Join("peripherals", "adc")
}

func (l ADC) WriteFiles() []generator.WriteFile {
	return []generator.WriteFile{
		{
			FileName:     "adc.c",
			TemplateName: "adc.c",
		},
		{
			FileName:     "adc.h",
			TemplateName: "adc.h",
		},
	}
}

func (l ADC) Includes() []string {
	return []string{"zephyr/drivers/adc.h", "adc.h"}
}

func (l ADC) ApplyOverlay(dt *devicetree.DeviceTree) error {
	for _, instance := range l.Instances {
		if err := instance.AttachSelf(dt); err != nil {
			return fmt.Errorf("attach adc: %w", err)
		}
	}

	return nil
}
