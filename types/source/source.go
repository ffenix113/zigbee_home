package source

import (
	"github.com/ffenix113/zigbee_home/config"
	"github.com/ffenix113/zigbee_home/templates"
	"github.com/ffenix113/zigbee_home/types/generator"
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
