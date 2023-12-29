package source

import (
	"github.com/ffenix113/zigbee_home/cli/config"
	"github.com/ffenix113/zigbee_home/cli/templates"
	"github.com/ffenix113/zigbee_home/cli/types/generator"
)

type Source struct {
	templates *templates.Templates
}

func NewSource() *Source {
	return &Source{
		templates: templates.NewTemplates(templates.TemplateFS),
	}
}

func (s *Source) WriteTo(srcDir string, device *config.Device, extenders []generator.Extender) error {
	return s.templates.WriteTo(srcDir, device, extenders)
}
