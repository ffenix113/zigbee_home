package yamlstrict

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

func Unmarshal(to any, node *yaml.Node) error {
	marshaled, err := yaml.Marshal(node)
	if err != nil {
		return fmt.Errorf("marshaling node to bytes: %w", err)
	}

	dec := yaml.NewDecoder(bytes.NewReader(marshaled))
	dec.KnownFields(true)

	err = dec.Decode(to)
	if err != nil {
		return fmt.Errorf("unmarshal strict: %w", err)
	}

	return nil
}
