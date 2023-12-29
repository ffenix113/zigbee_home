package generator

import (
	"fmt"
	"strings"

	"github.com/ffenix113/zigbee_home/cli/types/appconfig"
	"github.com/ffenix113/zigbee_home/cli/types/devicetree"
)

type Adder interface {
	AppConfig() []appconfig.ConfigValue
	ApplyOverlay(overlay *devicetree.DeviceTree) error
}

type WriteFile struct {
	FileName     string
	TemplateName string
}

// Extender provides a way to extend source code by writing new files
// and add code to main.c
//
// Each unique extender will be executed only once.
// Uniqueness of extender is determined either by package & type name,
// or by `Name`, if extender is instance of `SimpleExtender`.
type Extender interface {
	// Includes returns include paths that will be
	// added after all pre-defined includes in main.c
	Includes() []string
	// Template returns template path that will define
	// known extender steps.
	Template() string
	// WriteFiles returns a slice of files that should be written
	// from the template names. Necessary paths will be
	// created if file is located in directory that does not yet exist.
	// If any headers are created - they will not be included in main.c
	// If this is needed - add file path to `Includes()` return value.
	WriteFiles() []WriteFile
}

var _ Extender = SimpleExtender{}

type SimpleExtender struct {
	// Name is unique identifier of extender.
	// Used only for deduplication purpuses.
	Name           string
	IncludeHeaders []string
	TemplateName   string
	FilesToWrite   []WriteFile
}

func ExtenderName(e Extender) string {
	extenderName := strings.TrimPrefix(fmt.Sprintf("%T", e), "*")
	if simpleExtender, ok := e.(SimpleExtender); ok {
		extenderName = simpleExtender.Name
		if extenderName == "" {
			panic("all simple extenders require `Name` field to be set")
		}
	}

	return extenderName
}

// Includes implements Extender.
func (e SimpleExtender) Includes() []string {
	return e.IncludeHeaders
}

// Template implements Extender.
func (e SimpleExtender) Template() string {
	return e.TemplateName
}

// WriteFiles implements Extender.
func (e SimpleExtender) WriteFiles() []WriteFile {
	return e.FilesToWrite
}
