package types

import (
	"fmt"

	"github.com/ffenix113/zigbee_home/types/yamlstrict"
	"gopkg.in/yaml.v3"
)

var _ yaml.Unmarshaler = (*Option[byte])(nil)

type Option[T any] struct {
	hasValue bool
	value    T
}

func NewOption[T any](v T) Option[T] {
	return Option[T]{
		hasValue: true,
		value:    v,
	}
}

func NewEmptyOption[T any]() Option[T] {
	return Option[T]{}
}

func (o *Option[T]) Set(v T) {
	o.hasValue = true
	o.value = v
}

func (o *Option[T]) Reset() {
	(*o) = Option[T]{}
}

func (o Option[T]) HasValue() bool {
	return o.hasValue
}

func (o Option[T]) Value() T {
	return o.value
}

func (o Option[T]) String() string {
	return fmt.Sprint(o.Value())
}

// Format is a proxy for Option to format value using some verb.
func (o Option[T]) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, fmt.FormatString(f, verb), o.Value())
}

func (o *Option[T]) UnmarshalYAML(n *yaml.Node) error {
	var val T

	if err := yamlstrict.Unmarshal(&val, n); err != nil {
		return fmt.Errorf("unmarshal option with type %T: %w", val, err)
	}

	o.hasValue = true
	o.value = val

	return nil
}
