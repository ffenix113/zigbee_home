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
	})

	t.templates = template.Must(tpls.ParseFS(templateFS, "src/*.tpl", "src/**/*.tpl"))

	return t
}

func (t *Templates) WriteTo(srcDir string, device *types.TemplateDevice) error {
	file, err := os.Create(srcDir + "/main.c")
	if err != nil {
		return fmt.Errorf("create main.c: %w", err)
	}

	ctx := Context{
		GeneratedOn: time.Now().UTC(),
		Version:     "0.0.0-dev",
		Device:      device,
	}

	return t.templates.ExecuteTemplate(file, "main.c.tpl", ctx)
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

var knownClusterAttrTemplates = map[cluster.ID]string{
	cluster.ID_TEMP_MEASUREMENT: "define_temperature_attr_list",
}
