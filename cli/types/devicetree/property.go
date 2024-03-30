package devicetree

import (
	"fmt"
	"io"
	"strings"
)

var PropertyStatusEnable = NewProperty(PropertyNameStatus, StatusOkay)
var PropertyStatusDisable = NewProperty(PropertyNameStatus, StatusDisabled)

const StatusOkay = rawValue(`"okay"`)
const StatusDisabled = rawValue(`"disabled"`)

const PropertyNameCompatible = "compatible"
const PropertyNameStatus = "status"

type Property struct {
	Name  string
	Value PropertyValue
}

func NewProperty(name string, value PropertyValue) Property {
	return Property{name, value}
}

func (p Property) writeTo(w io.StringWriter) error {
	w.WriteString(p.Name)

	if p.Value != nil {
		w.WriteString(" = " + p.Value.Value())
	}

	w.WriteString(";")

	return nil
}

type PropertyValue interface {
	Value() string
}

type rawValue string

func (v rawValue) Value() string {
	return string(v)
}

func PropertyValueFn(fn func() string) PropertyValue {
	return rawValue(fn())
}

func String(value string) PropertyValue {
	return rawValue(value)
}

func Quoted(value string) PropertyValue {
	return rawValue(`"` + value + `"`)
}

func Label(label string) PropertyValue {
	return rawValue("&" + label)
}

// NrfPSel
// Reference: https://docs.zephyrproject.org/apidoc/latest/nrf-pinctrl_8h.html
func NrfPSel(fun string, port, pin uint8) PropertyValue {
	formatted := fmt.Sprintf("NRF_PSEL(%s, %d, %d)", fun, port, pin)

	return Angled(rawValue(formatted))
}

func Angled(value PropertyValue) PropertyValue {
	return rawValue("<" + value.Value() + ">")
}

func Array(values ...PropertyValue) PropertyValue {
	if len(values) < 2 {
		panic("array must have at least two values")
	}

	parts := make([]string, 0, len(values))

	for _, value := range values {
		parts = append(parts, value.Value())
	}

	return rawValue(strings.Join(parts, ", "))
}

func FromValue(val any) PropertyValue {
	switch typed := val.(type) {
	case string:
		return Quoted(typed)
	case uint8, int, int8:
		return Angled(rawValue(fmt.Sprintf("%d", val)))
	}

	panic(fmt.Sprintf("unknown type to convert to property value: %T", val))
}
