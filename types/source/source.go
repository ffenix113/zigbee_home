package source

import (
	"github.com/ffenix113/zigbee_home/config"
	"github.com/ffenix113/zigbee_home/templates"
	"github.com/ffenix113/zigbee_home/types"
	"github.com/ffenix113/zigbee_home/types/generator"
)

type Source struct {
	templates *templates.Templates
}

func NewSource(ncsVersion types.Semver, templatesPath string) *Source {
	return &Source{
		templates: templates.NewTemplates(templates.TemplateFS(templatesPath), ncsVersion),
	}
}

func (s *Source) WriteTo(srcDir string, device *config.Device, extenders []generator.Extender) error {
	return s.templates.WriteTo(srcDir, device, extenders)
}
