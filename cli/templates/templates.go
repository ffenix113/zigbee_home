package templates

import (
	"text/template"

	"github.com/ffenix113/zigbee_home/cli/config"
)

type Templates struct {
	parsed *template.Template
}

func (Templates) Prepare(cfg config.Device, templatesDir string) error {
	for sensorName := range cfg.Sensors {
		// Do something.
	}
}

func (Templates) WriteTo(srcDir string) error {

}
