package devicetree

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

const StatusOkay = rawValue(`"okay"`)
const StatusDisabled = rawValue(`"disabled"`)

const NodeNameRoot = "/"
const NodeNameChosen = "chosen"

const PropertyNameCompatible = "compatible"

type DeviceTree struct {
	Nodes []Node
}

func NewDeviceTree() *DeviceTree {
	return &DeviceTree{}
}

func (t *DeviceTree) WriteTo(w io.StringWriter) error {
	for _, node := range t.Nodes {
		if err := node.WriteTo(w); err != nil {
			return fmt.Errorf("write DeviceTree: %w", err)
		}
	}

	return nil
}

func (t *DeviceTree) AddNodes(nodes ...Node) *DeviceTree {
	t.Nodes = append(t.Nodes, nodes...)

	return t
}

type Node struct {
	Name, Label string
	UnitAddress string

	Properties []Property

	SubNodes []Node
}

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

const identationSymbol = "\t"

func (n *Node) WriteTo(w io.StringWriter) error {
	return n.writeTo("", w)
}

func (n *Node) writeTo(identation string, w io.StringWriter) error {
	w.WriteString(identation)
	if n.Label != "" {
		w.WriteString(n.Label + ": ")
	}
	w.WriteString(n.Name)

	if n.UnitAddress != "" {
		w.WriteString("@" + n.UnitAddress)
	}

	w.WriteString(" {\n")

	nodeIdent := identation + identationSymbol

	for _, property := range n.Properties {
		w.WriteString(nodeIdent)

		property.writeTo(w)

		w.WriteString("\n")
	}

	for _, subNode := range n.SubNodes {
		subNode.writeTo(nodeIdent, w)
	}

	w.WriteString(identation + "};\n")

	return nil
}

func (v rawValue) Value() string {
	return string(v)
}

func PropertyValueFn(fn func() string) PropertyValue {
	return rawValue(fn())
}

func String(value string) PropertyValue {
	return rawValue(`"` + value + `"`)
}

func Label(label string) PropertyValue {
	return rawValue("&" + label)
}

// NrfPSel
// Reference: https://docs.zephyrproject.org/apidoc/latest/nrf-pinctrl_8h.html
func NrfPSel(fun string, port, pin int) PropertyValue {
	portStr := strconv.Itoa(port)
	pinStr := strconv.Itoa(pin)

	return Angled("NRF_PSEL(" + fun + ", " + portStr + ", " + pinStr + ")")
}

func Angled(value string) PropertyValue {
	return rawValue("<" + value + ">")
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
		return String(typed)
	case int:
		return Angled(strconv.Itoa(typed))
	}

	panic(fmt.Sprintf("unknown type to convert to property value: %T", val))
}
