package templates

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/ffenix113/zigbee_home/config"
	"github.com/ffenix113/zigbee_home/types"
	"github.com/ffenix113/zigbee_home/types/generator"
	"github.com/ffenix113/zigbee_home/types/sensor"
	"github.com/ffenix113/zigbee_home/zcl/cluster"
)

// TemplateFS is for sensor templates.
// For example src/extenders/sensors/bosch/bme280.tpl
//
//go:embed src/*.hpp src/*.cpp src/*.tpl src/*/*.tpl src/*/*/*.tpl
//go:embed src/modules/*/dts/bindings/sensor/*.yaml src/modules/*/zephyr/*
var embTemplateFS embed.FS

var TemplateFS = func(templatesPath string) fs.FS {
	if templatesPath == "" {
		return embTemplateFS
	}

	// return embTemplateFS
	// FIXME: This is for debug, so should be at least put behind some option.
	expanded, err := filepath.Abs(os.ExpandEnv(templatesPath))
	if err != nil {
		panic(err)
	}

	return os.DirFS(expanded)
}

// This map can be removed in favor of cluster telling
// which template it want's to use, or try
// CVarName value as template as well.
var knownClusterTemplates = map[cluster.ID]string{
	cluster.ID_BASIC:                     "basic",
	cluster.ID_POWER_CONFIG:              "power_config",
	cluster.ID_DEVICE_TEMP_CONFIG:        "device_temp_config",
	cluster.ID_ON_OFF:                    "on_off",
	cluster.ID_TEMP_MEASUREMENT:          "temperature",
	cluster.ID_REL_HUMIDITY_MEASUREMENT:  "water_content",
	cluster.ID_PRESSURE_MEASUREMENT:      "pressure",
	cluster.ID_CARBON_DIOXIDE:            "carbon_dioxide",
	cluster.ID_IAS_ZONE:                  "ias_zone",
	cluster.ID_SOIL_MOISTURE_MEASUREMENT: "water_content",
}

var sourceFiles = [][2]string{
	{path.Join("..", "CMakeLists.txt"), "CMakeLists.txt.tpl"},
	{path.Join("..", "Kconfig"), "Kconfig.tpl"},
	{path.Join("..", "Kconfig.sysbuild"), "Kconfig.sysbuild.tpl"},
	{"main.cpp", "main.cpp.tpl"},
	{"device.hpp", "device.hpp.tpl"},
	{"clusters.hpp", "clusters.hpp.tpl"},
	{"types.hpp", "types.hpp"},
	{"types.cpp", "types.cpp"},
	{"types_button_handler.hpp", "types_button_handler.hpp"},
	{"types_button_handler.cpp", "types_button_handler.cpp"},
}

var knownExtenders = [...]string{
	"include",
	"top_level",
	"main",
	"attr_init",
	"loop",
}

type Templates struct {
	templates    *template.Template
	templateTree templateTree
}

type SensorCtx struct {
	Endpoint int
	Sensor   sensor.Sensor
	Device   *config.Device
	Extender generator.Extender
}

type ClusterCtx struct {
	Endpoint int
	Cluster  cluster.Cluster
}

type Context struct {
	GeneratedOn time.Time
	Version     string
	// Each sensor will have it's own endpoint,
	// which should be okay for now.
	// Packing multiple sensor into one endpoint
	// can be done in the future.
	Device *config.Device

	Extenders []generator.Extender
}

type ContextWithAdditional struct {
	Context
	Extender          generator.Extender
	Sensor            sensor.Sensor
	AdditionalContext any
}

func NewTemplates(templateFS fs.FS, ncsVersion types.Semver) *Templates {
	t := &Templates{
		templateTree: templateTree{
			tree: make(map[string]*templateTree),
		},
		templates: template.Must(template.New("").Parse("")),
	}

	t.templates.Funcs(template.FuncMap{
		"typeFromSensor":      typeFromSensor,
		"clusterTpl":          t.clusterTpl,
		"render":              t.render,
		"maybeRender":         t.maybeRender,
		"maybeRenderExtender": t.maybeRenderExtender,
		"toButtonBit":         toButtonBit,
		"sensorCtx":           sensorCtx,
		"clusterCtx":          clusterCtx,
		"isLast":              isLast,
		"sum":                 sum,
		"formatHex":           formatHex,
		"joinPath": func(strs ...string) string {
			return path.Join(strs...)
		},
		// Specific functions to check exact version
		// so we would know where each one is used,
		// and what we can deprecate.
		"ncsVersionIs_2_6": ncsVersionIs(ncsVersion, types.Semver{2, 6, 0}),
	})

	must(t.parseByDir(templateFS, nil))

	t.templates = template.Must(t.templates.ParseFS(templateFS,
		path.Join("src", "*.tpl"),

		path.Join("src", "*.cpp"),
		path.Join("src", "*.hpp"),

		path.Join("src", "zigbee", "*.tpl")),
	)

	return t
}

func (t *Templates) parseByDir(tplFS fs.FS, validateTpl func(t *template.Template) error) error {
	// FIXME: it

	err := fs.WalkDir(tplFS, "src", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		tplFile := path

		openTpl, err := tplFS.Open(tplFile)
		if err != nil {
			return fmt.Errorf("open template %q: %w", tplFile, err)
		}
		defer openTpl.Close()

		newTpl := templateFromPath(&t.templateTree, t.templates, strings.TrimPrefix(tplFile, "src/"))

		tplText, err := io.ReadAll(openTpl)
		if err != nil {
			return fmt.Errorf("read template %q: %w", tplFile, err)
		}

		if _, err := newTpl.Parse(string(tplText)); err != nil {
			return fmt.Errorf("parse template %q: %w", tplFile, err)
		}

		if validateTpl != nil {
			if err := validateTpl(newTpl); err != nil {
				return fmt.Errorf("validate template %q: %w", tplFile, err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("parse template fs: %w", err)
	}

	return nil
}

func templateFromPath(root *templateTree, baseTpl *template.Template, tplPath string) *template.Template {
	const templateExtention = ".tpl"

	pathParts := strings.Split(tplPath, "/") // because we always need to use "/" when using embed.FS

	tree := root
	for _, pathPart := range pathParts[:len(pathParts)-1] {
		if tree.tree == nil {
			tree.tree = make(map[string]*templateTree, 1)
		}

		subTree, ok := tree.tree[pathPart]
		if !ok {
			subTree = &templateTree{
				tree: make(map[string]*templateTree),
			}
			tree.tree[pathPart] = subTree
		}

		tree = subTree
	}

	templateName := pathParts[len(pathParts)-1]
	// FIXME: Removing the suffix results in inconsistent template names.
	pathPart, _ := strings.CutSuffix(templateName, templateExtention)

	subTree, ok := tree.tree[pathPart]
	if !ok {
		subTree = &templateTree{
			tree: make(map[string]*templateTree),
		}
		tree.tree[pathPart] = subTree
	}

	tree = subTree

	tree.tpl, _ = baseTpl.Clone()

	tplName, _ := strings.CutSuffix(tplPath, templateExtention)
	tree.tpl = tree.tpl.New(tplName)

	return tree.tpl
}

func (t *Templates) WriteTo(srcDir string, device *config.Device, extenders []generator.Extender) error {
	ctx := Context{
		GeneratedOn: time.Now().UTC(),
		Version:     "0.0.0-dev",
		Device:      device,

		Extenders: extenders,
	}

	// It is possible that sensors or extenders would write duplicate files.
	// Currently it should be okay, but it should be looked out for.
	// Maybe one of those files will contain templated values that would be overwritten..

	for _, sourceDefinition := range sourceFiles {
		template := t.templates.Lookup(sourceDefinition[1])
		if template == nil {
			return fmt.Errorf("tried to write unknown template: %q", sourceDefinition[1])
		}

		if err := writeTemplate(template, path.Join(srcDir, sourceDefinition[0]), ctx); err != nil {
			return fmt.Errorf("write template: %w", err)
		}
	}

	for _, extender := range extenders {
		if err := t.verifyExtender(extender); err != nil {
			return fmt.Errorf("extender %q is invalid: %w", generator.ExtenderName(extender), err)
		}

		// Files required by extender. Could be some implementation or helper functions.
		for _, fileToWrite := range extender.WriteFiles() {
			if fileToWrite.FileName == "" {
				fileToWrite.FileName = fileToWrite.TemplateName
			}

			template := t.findTemplate(fileToWrite.TemplateName)
			if err := writeTemplate(
				template,
				filepath.Join(srcDir, fileToWrite.FileName),
				ContextWithAdditional{Context: ctx, Extender: extender, AdditionalContext: fileToWrite.AdditionalContext}); err != nil {
				return fmt.Errorf("write extender file %q: %w", fileToWrite.FileName, err)
			}
		}

		// Module files. For example drivers for sensors.
		for _, module := range extender.ZephyrModules() {
			moduleTemplates := (&t.templateTree).FindByPath("modules", module)
			for _, moduleTemplate := range moduleTemplates {
				// Modules are not in 'src' dir, but in the root instead.
				if err := writeTemplate(moduleTemplate, filepath.Join(srcDir, "..", moduleTemplate.Name()), ctx); err != nil {
					return fmt.Errorf("write extender module %q file %q: %w", module, moduleTemplate.Name(), err)
				}
			}
		}
	}

	for _, sensor := range device.Sensors {
		fileWriter, ok := sensor.(interface {
			WriteFiles() []generator.WriteFile
		})
		if !ok {
			continue
		}

		filesToWrite := fileWriter.WriteFiles()

		for _, fileToWrite := range filesToWrite {
			if fileToWrite.FileName == "" {
				fileToWrite.FileName = fileToWrite.TemplateName
			}

			template := t.findTemplate(fileToWrite.TemplateName)
			if err := writeTemplate(
				template,
				filepath.Join(srcDir, fileToWrite.FileName),
				ContextWithAdditional{Context: ctx, Sensor: sensor, AdditionalContext: fileToWrite.AdditionalContext}); err != nil {
				return fmt.Errorf("write sensor file %q: %w", fileToWrite.FileName, err)
			}
		}
	}

	return nil
}

func (t *Templates) verifyExtender(extender generator.Extender) error {
	if extender.Template() == "" {
		return nil
	}

	tpl := t.findTemplate(extender.Template())
	if tpl == nil {
		return fmt.Errorf("required extention template not found: %q", extender.Template())
	}

	var foundExtenders int

	for _, knownExtender := range knownExtenders {
		extdTpl := tpl.Lookup(knownExtender)
		if extdTpl != nil {
			foundExtenders++
		}
	}

	if foundExtenders != len(tpl.Templates())-2 {
		return fmt.Errorf("extender template has weird templates: %s", tpl.DefinedTemplates())
	}

	return nil
}

func writeTemplate(template *template.Template, filePath string, ctx any) error {
	if template == nil {
		return fmt.Errorf("template is nil for path %q", filePath)
	}

	dirName := filepath.Dir(filePath)
	if err := os.MkdirAll(dirName, 0o755); err != nil {
		return fmt.Errorf("create directory %q: %w", dirName, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create %q: %w", filePath, err)
	}

	defer file.Close()

	if err := template.Execute(file, ctx); err != nil {
		return fmt.Errorf("execute source template %q: %w", template.Name(), err)
	}

	return nil
}

func (t *Templates) findTemplate(templateName string) *template.Template {
	tpl := t.templates.Lookup(templateName)
	if tpl != nil {
		return tpl
	}

	possiblePaths := [][]string{
		strings.Split(templateName, "/"),
		append([]string{"extenders"}, strings.Split(templateName, "/")...), // because we always need to use "/" when using embed.FS
	}

	for _, possiblePath := range possiblePaths {
		tree := &t.templateTree

		for _, namePart := range possiblePath {
			tree = tree.tree[namePart]
			if tree == nil {
				break
			}
		}

		if tree == nil || tree.tpl == nil {
			continue
		}

		return tree.tpl
	}

	return nil
}

func (t *Templates) clusterTpl(clusterID cluster.ID, tplSuffix string) (string, error) {
	tplName, ok := knownClusterTemplates[clusterID]
	if !ok {
		val, _ := clusterID.ToZCL()
		return "", fmt.Errorf("unknown cluster ID: %q(%d)", val, clusterID)
	}

	return tplName + "_" + tplSuffix, nil
}

func (t *Templates) render(tplName string, ctx any) (string, error) {
	var buf bytes.Buffer

	if err := t.templates.ExecuteTemplate(&buf, tplName, ctx); err != nil {
		return "", fmt.Errorf("execute %q: %w", tplName, err)
	}

	return buf.String(), nil
}

// maybeRender will render template `tplName`, if it exists.
// This is useful for optional templates.
func (t *Templates) maybeRender(tplName string, ctx any) (string, error) {

	tpl := t.templates.Lookup(tplName)
	if tpl == nil {
		return "", nil
	}

	var buf bytes.Buffer
	if err := t.templates.ExecuteTemplate(&buf, tplName, ctx); err != nil {
		return "", fmt.Errorf("execute %q: %w", tplName, err)
	}

	return buf.String(), nil
}

// maybeRender will render template `tplName`, if it exists.
// This is useful for optional templates.
func (t *Templates) maybeRenderExtender(tplPath, tplName string, ctx any) (string, error) {
	if tplPath == "" {
		return "", nil
	}

	tpl := t.findTemplate(tplPath)
	if tpl == nil {
		return "", fmt.Errorf("extender template %q is not defined", tplPath)
	}

	tpl = tpl.Lookup(tplName)
	if tpl == nil {
		return "", nil
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, ctx); err != nil {
		return "", fmt.Errorf("execute %q -> %q: %w", tplPath, tplName, err)
	}

	return buf.String(), nil
}

func typeFromSensor(sensor SensorCtx) string {
	return sensor.Sensor.CPPComponentType()
}

func toButtonBit(btnID string) string {
	return "BUTTON_BIT(" + btnID + ")"
}

func sensorCtx(endpoint int, device *config.Device, sensor sensor.Sensor, extender generator.Extender) SensorCtx {
	return SensorCtx{
		Endpoint: endpoint,
		Device:   device,
		Sensor:   sensor,
		Extender: extender,
	}
}

func clusterCtx(endpoint int, cluster cluster.Cluster) ClusterCtx {
	return ClusterCtx{
		Endpoint: endpoint,
		Cluster:  cluster,
	}
}

func isLast(i, arrLen int) bool {
	return i+1 == arrLen
}

func sum(a, b int) int {
	return a + b
}

func formatHex(val any) (string, error) {
	switch i := val.(type) {
	case uint8, uint16, uint32, uint64, uint,
		int8, int16, int32, int64, int:
		return fmt.Sprintf("%#x", i), nil
	default:
		return "", fmt.Errorf("unknown type to format: %T", val)
	}
}

func ncsVersionIs(current, another types.Semver) func() bool {
	isSame := current.SameMajorMinor(another)

	return func() bool {
		return isSame
	}
}

func must(err error) {
	if err == nil {
		return
	}

	errText := err.Error()

	if strings.Contains(errText, "template: pattern matches no") {
		return
	}

	panic(err)
}
