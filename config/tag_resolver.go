package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// tagsResolver does required actions to resolve tags
// inside configuration file.
// For example does inclusion of external files.
type tagsResolver struct {
	// maxIncludeDepth tells how deep includes can be
	// when file includes file, includes file, ...
	maxIncludeDepth int
}

func newTagsResolver() *tagsResolver {
	const defaultMaxIncludeDepth = 3

	return &tagsResolver{
		maxIncludeDepth: defaultMaxIncludeDepth,
	}
}

func (r *tagsResolver) resolve(node *yaml.Node, depth int) error {
	var newNode *yaml.Node

	var err error

	switch node.Tag {
	case "!include":
		newNode, err = r.getIncludedNode(node.Value)
		if err != nil {
			return fmt.Errorf("include: %w", err)
		}

		if depth+1 > r.maxIncludeDepth {
			log.Fatalf("max include resolution depth of %d reached", depth)
		}

		if err := r.resolve(newNode, depth+1); err != nil {
			return fmt.Errorf("included file: %w", err)
		}
	case "!env":
		newNode, err = r.getEnvNode(node.Value)
		if err != nil {
			return fmt.Errorf("env: %w", err)
		}
	}

	if newNode != nil {
		*node = *newNode
	}

	for _, content := range node.Content {
		if err := r.resolve(content, depth); err != nil {
			return fmt.Errorf("content: %w", err)
		}
	}

	return nil
}

func (r *tagsResolver) getIncludedNode(includePath string) (*yaml.Node, error) {
	includePath, err := filepath.Abs(includePath)
	if err != nil {
		return nil, fmt.Errorf("get absolute include path of %q: %w", includePath, err)
	}

	fileBts, err := os.ReadFile(includePath)
	if err != nil {
		return nil, fmt.Errorf("read file %q: %w", includePath, err)
	}

	var n yaml.Node
	if err := yaml.Unmarshal(fileBts, &n); err != nil {
		return nil, fmt.Errorf("unmarshal %q: %w", includePath, err)
	}

	return r.getUnmarshaledNode(&n), nil
}

func (r *tagsResolver) getEnvNode(env string) (*yaml.Node, error) {
	envValue, ok := os.LookupEnv(env)
	if !ok || envValue == "" {
		return nil, nil
	}

	var n yaml.Node
	if err := yaml.Unmarshal([]byte(envValue), &n); err != nil {
		return nil, fmt.Errorf("unmarshal %q: %w", envValue, err)
	}

	return r.getUnmarshaledNode(&n), nil
}

func (r *tagsResolver) getUnmarshaledNode(n *yaml.Node) *yaml.Node {
	if n.Kind == yaml.DocumentNode {
		return n.Content[0]
	}

	return n
}
