package generator

import (
	"fmt"
	"strings"

	"github.com/ffenix113/zigbee_home/types/appconfig"
	"github.com/ffenix113/zigbee_home/types/devicetree"
)

type Adder interface {
	AppConfig() []appconfig.ConfigValue
	ApplyOverlay(overlay *devicetree.DeviceTree) error
}

type WriteFile struct {
	FileName          string
	TemplateName      string
	AdditionalContext any
}

// Extender provides a way to extend source code by writing new
// files, add code to main.c and extend CMakeLists.txt
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
	// ZephyrModules adds a list of modules that extender provides.
	// Module can be anything(I guess).
	//
	// Module name will also write all files in templates
	// from the `modules/<zephyrModule>` path for each provided module.
	//
	// https://docs.zephyrproject.org/latest/develop/modules.html
	ZephyrModules() []string
}

var _ Adder = SimpleExtender{}
var _ Extender = SimpleExtender{}

type SimpleExtender struct {
	// Name is unique identifier of extender.
	// Used only for deduplication purpuses.
	Name              string
	IncludeHeaders    []string
	TemplateName      string
	FilesToWrite      []WriteFile
	ZephyrModuleNames []string
	Config            []appconfig.ConfigValue
	OverlayFn         func(overlay *devicetree.DeviceTree) error
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

func (e SimpleExtender) ZephyrModules() []string {
	return e.ZephyrModuleNames
}

func (e SimpleExtender) AppConfig() []appconfig.ConfigValue {
	return e.Config
}

func (e SimpleExtender) ApplyOverlay(overlay *devicetree.DeviceTree) error {
	if e.OverlayFn == nil {
		return nil
	}

	return e.OverlayFn(overlay)
}
