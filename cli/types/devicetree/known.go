package devicetree

import (
	"github.com/ffenix113/zigbee_home/cli/types"
)

type KnownNode interface {
	AttachSelf(dt *DeviceTree) error
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
						NewProperty("psels", NrfPSel("UART_TX", u.Tx.Port, u.Tx.Pin)),
					},
				},
				{
					Name: "group2",
					Properties: []Property{
						NewProperty("psels", NrfPSel("UART_RX", u.Rx.Port, u.Rx.Pin)),
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
						NrfPSel("UART_TX", u.Tx.Port, u.Tx.Pin),
						NrfPSel("UART_RX", u.Rx.Port, u.Rx.Pin),
					)),
					NewProperty("low-power-enable", nil),
				},
			},
		},
	})

	return nil
}
