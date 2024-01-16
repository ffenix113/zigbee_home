package templates

import (
	"fmt"
	"strings"
	"text/template"
)

type templateTree struct {
	tree map[string]*templateTree
	tpl  *template.Template
}

func (t *templateTree) FindByPath(prefix string) []*template.Template {
	parts := strings.Split(prefix, "/")

	currentTree := t
	for _, part := range parts {
		currentTree = currentTree.tree[part]
		if currentTree == nil {
			panic(fmt.Sprintf("part %q of prefix %q is not present in template tree", part, prefix))
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
