package templates

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"text/template"
	"time"

	"github.com/ffenix113/zigbee_home/cli/types"
	"github.com/ffenix113/zigbee_home/cli/zcl/cluster"
)

//go:embed src/*.tpl src/**/*.tpl
var TemplateFS embed.FS

type Templates struct {
	templates *template.Template
}

type Context struct {
	GeneratedOn time.Time
	Version     string
	Device      *types.TemplateDevice
}

func NewTemplates(templateFS fs.FS) *Templates {
	t := &Templates{}

	tpls := template.New("")
	tpls.Funcs(template.FuncMap{
		"clusterToAttrTemplate": t.clusterToAttrTemplate,
		"render":                t.render,
		"isLast":                isLast,
		"sum":                   sum,
	})

	t.templates = template.Must(tpls.ParseFS(templateFS, "src/*.tpl", "src/**/*.tpl"))

	return t
}

func (t *Templates) WriteTo(srcDir string, device *types.TemplateDevice) error {
	ctx := Context{
		GeneratedOn: time.Now().UTC(),
		Version:     "0.0.0-dev",
		Device:      device,
	}

	for _, sourceDefinition := range sourceFiles {
		file, err := os.Create(srcDir + "/" + sourceDefinition[0])
		if err != nil {
			return fmt.Errorf("create %q: %w", sourceDefinition[0], err)
		}
		// This will be executed when `WriteTo` exits,
		// but it is fine for us.
		defer file.Close()

		if err := t.templates.ExecuteTemplate(file, sourceDefinition[1], ctx); err != nil {
			return fmt.Errorf("execute source template %q: %w", sourceDefinition[1], err)
		}
	}

	return nil
}

var sourceFiles = [][2]string{
	{"../CMakeLists.txt", "CMakeLists.txt.tpl"},
	{"main.c", "main.c.tpl"},
	{"device.h", "device.h.tpl"},
	{"clusters.h", "clusters.h.tpl"},
}

func (t *Templates) clusterToAttrTemplate(cluster cluster.Cluster) (string, error) {
	tplName, ok := knownClusterAttrTemplates[cluster.ID()]
	if !ok {
		return "", fmt.Errorf("unknown cluster ID: %q", cluster.ID().String())
	}

	return tplName, nil
}

func (t *Templates) render(tplName string, ctx any) (string, error) {
	var buf bytes.Buffer

	if err := t.templates.ExecuteTemplate(&buf, tplName, ctx); err != nil {
		return "", fmt.Errorf("execute %q: %w", tplName, err)
	}

	return buf.String(), nil
}

func isLast(i, arrLen int) bool {
	return i+1 == arrLen
}

func sum(a, b int) int {
	return a + b
}

var knownClusterAttrTemplates = map[cluster.ID]string{
	cluster.ID_TEMP_MEASUREMENT:         "define_temperature_attr_list",
	cluster.ID_REL_HUMIDITY_MEASUREMENT: "define_humidity_attr_list",
	cluster.ID_PRESSURE_MEASUREMENT:     "define_pressure_attr_list",
	cluster.ID_ON_OFF:                   "define_on_off_attr_list",
}
