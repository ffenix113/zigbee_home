package devicetree

import (
	"fmt"
	"strings"

	"github.com/ffenix113/zigbee_home/cli/types"
)

type KnownNode interface {
	AttachSelf(dt *DeviceTree) error
}

// Verify that pin implements both interfaces
var _ interface {
	LED
	Button
} = pin{}

type LED interface {
	KnownNode
	led()
}

type Button interface {
	KnownNode
	button()
}

type UART struct {
	Tx, Rx types.Pin
}

func (u UART) AttachSelf(dt *DeviceTree) error {
	pinctrlNode := dt.FindSpecificNode(SearchByLabel(NodeLabelPinctrl))
	if pinctrlNode == nil {
		return ErrNodeNotFound(NodeLabelPinctrl)
	}

	pinctrlNode.AddNodes(
		&Node{
			Name:  "uart0_default",
			Label: "uart0_default",
			SubNodes: []*Node{
				{
					Name: "group1",
					Properties: []Property{
						NewProperty("psels", NrfPSel("UART_TX", u.Tx.Port.Value(), u.Tx.Pin.Value())),
					},
				},
				{
					Name: "group2",
					Properties: []Property{
						NewProperty("psels", NrfPSel("UART_RX", u.Rx.Port.Value(), u.Rx.Pin.Value())),
						NewProperty("bias-pull-up", nil),
					},
				},
			},
		},
	)

	pinctrlNode.AddNodes(&Node{
		Name:  "uart0_sleep",
		Label: "uart0_sleep",
		SubNodes: []*Node{
			{
				Name: "group1",
				Properties: []Property{
					NewProperty("psels", Array(
						NrfPSel("UART_TX", u.Tx.Port.Value(), u.Tx.Pin.Value()),
						NrfPSel("UART_RX", u.Rx.Port.Value(), u.Rx.Pin.Value()),
					)),
					NewProperty("low-power-enable", nil),
				},
			},
		},
	})

	return nil
}

func NewLED(ledPin types.Pin) LED {
	return pin{
		Pin: ledPin,

		nodeName:   "leds",
		compatible: "gpio-leds",
	}
}

func NewButton(btnPin types.Pin) Button {
	return pin{
		Pin: btnPin,

		nodeName:   "buttons",
		compatible: "gpio-keys",
	}
}

type pin struct {
	Pin types.Pin

	nodeName, compatible string
}

func (p pin) AttachSelf(dt *DeviceTree) error {
	pinName := p.Pin.ID
	if pinName == "" {
		pinName = p.Pin.Label()
	}

	aliases := dt.FindSpecificNode(SearchByName(NodeNameRoot), SearchByName(NodeNameAliases))
	aliases.Properties = append(aliases.Properties,
		NewProperty(strings.ReplaceAll(pinName, "_", "-"), Label(pinName)))

	// If pin is not defined - do not add its configuration.
	if !p.Pin.PinsDefined() {
		return nil
	}

	pinsNode := dt.FindSpecificNode(SearchByName(NodeNameRoot), SearchByName(p.nodeName))
	if pinsNode == nil {
		pinsNode = &Node{
			Name:  p.nodeName,
			Label: p.nodeName,
			Properties: []Property{
				NewProperty(PropertyNameCompatible, FromValue(p.compatible)),
			},
		}

		dt.FindSpecificNode(SearchByName(NodeNameRoot)).AddNodes(pinsNode)
	}

	activeState := "GPIO_ACTIVE_HIGH"
	if p.Pin.Inverted {
		activeState = "GPIO_ACTIVE_LOW"
	}

	pinsNode.AddNodes(&Node{
		Name:  pinName,
		Label: pinName,
		Properties: []Property{
			NewProperty("gpios", Angled(Label(fmt.Sprintf("gpio%d %d %s", p.Pin.Port, p.Pin.Pin, activeState)))),
		},
	})

	return nil
}

func (pin) led()    {}
func (pin) button() {}
