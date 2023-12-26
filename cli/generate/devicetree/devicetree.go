package devicetree

import (
	"fmt"
	"io"
)

const NodeNameRoot = "/"
const NodeNameChosen = "chosen"

const NodeLabelPinctrl = "pinctrl"
const NodeLabelI2c1 = "i2c1"

type ErrNodeNotFound string

func (f ErrNodeNotFound) Error() string {
	return fmt.Sprintf("node %q was not found", f)
}

type NodeSearchFn func(n *Node) bool

type DeviceTree struct {
	Nodes []*Node
}

type Node struct {
	Name, Label string
	UnitAddress string
	Upsert      bool

	Properties []Property

	SubNodes []*Node
}

func NewDeviceTree() *DeviceTree {
	return (&DeviceTree{}).
		AddNodes((&Node{Name: NodeNameRoot}).
			AddNodes(&Node{Name: NodeNameChosen})).
		AddNodes(&Node{Label: NodeLabelPinctrl, Upsert: true}).
		AddNodes(
			&Node{
				Label:  NodeLabelI2c1,
				Upsert: true,
				Properties: []Property{
					PropertyStatusEnable,
				},
			},
		)
}

func (t *DeviceTree) WriteTo(w io.StringWriter) error {
	for _, node := range t.Nodes {
		if err := node.WriteTo(w); err != nil {
			return fmt.Errorf("write DeviceTree: %w", err)
		}
	}

	return nil
}

func (t *DeviceTree) AddNodes(nodes ...*Node) *DeviceTree {
	t.Nodes = append(t.Nodes, nodes...)

	return t
}

func (t *DeviceTree) FindSpecificNode(searchFns ...NodeSearchFn) *Node {
	return findNode(t.Nodes, searchFns...)
}

func findNode(nodes []*Node, searchFns ...NodeSearchFn) *Node {
	if len(nodes) == 0 || len(searchFns) == 0 {
		return nil
	}

	for _, node := range nodes {
		if searchFns[0](node) {
			if len(searchFns) == 1 {
				return node
			}

			return findNode(node.SubNodes, searchFns[1:]...)
		}
	}

	return nil
}

const identationSymbol = "\t"

func (n *Node) Ref() string {
	return "&" + n.Label
}

func (n *Node) AddNodes(nodes ...*Node) *Node {
	n.SubNodes = append(n.SubNodes, nodes...)

	return n
}

func (n *Node) FindSpecificNode(searchFns ...NodeSearchFn) *Node {
	return findNode(n.SubNodes, searchFns...)
}

func (n *Node) WriteTo(w io.StringWriter) error {
	return n.writeTo("", w)
}

func (n *Node) writeTo(identation string, w io.StringWriter) error {
	if n.Upsert && n.Name != "" {
		return fmt.Errorf("node %q must not have name, as it is upserted", n.Name)
	}

	if !n.Upsert && n.Label != "" && n.Name == "" {
		return fmt.Errorf("node with label %q must have a name", n.Label)
	}

	w.WriteString(identation)
	if n.Label != "" {
		if n.Upsert {
			w.WriteString("&")
		}
		w.WriteString(n.Label)
		if !n.Upsert {
			w.WriteString(": ")
		}
	}
	if !n.Upsert {
		w.WriteString(n.Name)
	}

	if n.UnitAddress != "" {
		w.WriteString("@" + n.UnitAddress)
	}

	w.WriteString(" {\n")

	nodeIdent := identation + identationSymbol

	for _, property := range n.Properties {
		w.WriteString(nodeIdent)

		if err := property.writeTo(w); err != nil {
			return fmt.Errorf("write property %q to node %q: %w", property.Name, n.Ref(), err)
		}

		w.WriteString("\n")
	}

	for _, subNode := range n.SubNodes {
		if err := subNode.writeTo(nodeIdent, w); err != nil {
			return fmt.Errorf("write subNode %q to node %q: %w", subNode.Ref(), n.Ref(), err)
		}
	}

	w.WriteString(identation + "};\n\n")

	return nil
}

func SearchByLabel(label string) NodeSearchFn {
	return func(n *Node) bool {
		return n.Label == label
	}
}

func SearchByName(name string) NodeSearchFn {
	return func(n *Node) bool {
		return n.Name == name
	}
}
