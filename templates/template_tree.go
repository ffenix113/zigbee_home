package templates

import (
	"fmt"
	"path/filepath"
	"text/template"
)

type templateTree struct {
	tree map[string]*templateTree
	tpl  *template.Template
}

func (t *templateTree) FindByPath(parts ...string) []*template.Template {
	currentTree := t
	for _, part := range parts {
		currentTree = currentTree.tree[part]
		if currentTree == nil {
			panic(fmt.Sprintf("part %q of prefix %q is not present in template tree", part, filepath.Join(parts...)))
		}
	}

	return fetchTemplatesFromTree(currentTree, nil)
}

func fetchTemplatesFromTree(tree *templateTree, templates []*template.Template) []*template.Template {
	for _, innerTree := range tree.tree {
		templates = fetchTemplatesFromTree(innerTree, templates)
	}

	if tree.tpl != nil {
		templates = append(templates, tree.tpl)
	}

	return templates
}
