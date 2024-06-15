package types

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/ffenix113/zigbee_home/types/yamlstrict"
	"gopkg.in/yaml.v3"
)

var _ yaml.Unmarshaler = (*Pin)(nil)

// PinWithID is similar to Pin,
// but the difference is that user can define ID & short Pin for this type.
//
// So insted of `{id: 'some', port: 0, pin: 1}`
// user can provide `{id: 'some', pin: 0.04}`.
//
// Still, probably there is a more elegant way to achieve this.
type PinWithID struct {
	ID  string
	Pin Pin
}

func (p PinWithID) ToPin() Pin {
	if p.ID != "" && p.Pin.ID != "" {
		// This log message can be improved
		log.Fatalf("cannot have id & pin.id set at the same time for pin %#v", p)
	}

	if p.Pin.ID == "" {
		p.Pin.ID = p.ID
	}

	return p.Pin
}

type PinWithIDSlice []PinWithID

func (s PinWithIDSlice) ToPins() []Pin {
	pins := make([]Pin, 0, len(s))

	for _, pinWithID := range s {
		pins = append(pins, pinWithID.ToPin())
	}

	return pins
}

type Pin struct {
	ID       string
	Port     Option[uint8]
	Pin      Option[uint8]
	Inverted bool
}

func (p Pin) Name() string {
	if p.ID != "" {
		return p.ID
	}

	return p.Label()
}

func (p Pin) Label() string {
	return fmt.Sprintf("pin%s", p.NumericLabel())
}

func (p Pin) NumericLabel() string {
	return fmt.Sprintf("%d%d", p.Port, p.Pin)
}

func (p Pin) PinsDefined() bool {
	return p.Port.HasValue() && p.Pin.HasValue()
}

func (p Pin) Valid() bool {
	port := p.Port.Value()
	pin := p.Pin.Value()

	return p.ID != "" ||
		(p.PinsDefined() && ((port == 0 && pin <= 31) ||
			(port == 1 && pin <= 15)))
}

var pinRegex = regexp.MustCompile(`^([01])\.([0-3][0-9])$`)

func (p *Pin) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.ScalarNode {
		// A trick to remove custom unmarshaling behavior.
		type pin Pin
		if err := yamlstrict.Unmarshal((*pin)(p), value); err != nil {
			return fmt.Errorf("unmarshal pin: %w", err)
		}

		if !p.Valid() {
			return fmt.Errorf("pin has invalid definition(no pins & no id)")
		}

		return nil
	}

	if value.Value == "" {
		return fmt.Errorf("pin definition cannot be empty")
	}

	matches := pinRegex.FindStringSubmatch(value.Value)
	if matches == nil {
		return fmt.Errorf("pin definition must be in a form of X.XX, where X is a number")
	}

	port, err := strconv.ParseUint(matches[1], 10, 8)
	if err != nil {
		return fmt.Errorf("pin's port is invalid: %w", err)
	}

	pin, err := strconv.ParseUint(matches[2], 10, 8)
	if err != nil {
		return fmt.Errorf("pin's pin is invalid: %w", err)
	}

	p.Port = NewOption(uint8(port))
	p.Pin = NewOption(uint8(pin))

	return nil
}
