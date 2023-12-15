package templates

import (
	"io/fs"
	"text/template"

	"github.com/ffenix113/zigbee_home/cli/config"
)

type Templates struct {
	// parsed templates from templateFS
	parsed *template.Template
}

func (Templates) Prepare(cfg config.Device, templateFS fs.FS) error {
	for _, sensor := range cfg.Sensors {
		// Do something.
		_ = sensor
	}

	return nil
}

func (Templates) WriteTo(srcDir string) error {
	return nil
}
