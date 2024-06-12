package extenders

import (
	"strings"

	"github.com/ffenix113/zigbee_home/types/appconfig"
	"github.com/ffenix113/zigbee_home/types/devicetree"
	"github.com/ffenix113/zigbee_home/types/generator"
)

const DebugConsoleUSB = "usb"

type DebugLEDs struct {
	Enabled    bool
	Power      string
	Connection string
}

type DebugConfig struct {
	Enabled bool
	LEDs    DebugLEDs
	Console string
}

// Debug log optimizations
// CONFIG_DEBUG_OPTIMIZATIONS=y
// CONFIG_DEBUG_THREAD_INFO=y

func NewDebugUARTLog(config DebugConfig) generator.Extender {
	if !config.Enabled {
		return nil
	}

	if config.Console == "" {
		panic("console value must be set for debug configuration")
	}

	if strings.HasPrefix(config.Console, DebugConsoleUSB) && config.Console != DebugConsoleUSB {
		panic("debug backend should be 'usb', but is " + config.Console)
	}

	ledsEnabled := appconfig.Yes
	if !config.LEDs.Enabled {
		ledsEnabled = appconfig.No
	}

	return generator.SimpleExtender{
		IncludeHeaders: []string{
			"zephyr/logging/log.h",
			"zephyr/drivers/uart.h",
		},
		Config: []appconfig.ConfigValue{
			// Logging setup
			appconfig.CONFIG_LOG.Required(appconfig.Yes),
			appconfig.CONFIG_CONSOLE.Required(appconfig.Yes),
			appconfig.CONFIG_SERIAL.Required(appconfig.Yes),
			appconfig.CONFIG_LOG_BACKEND_UART.Required(appconfig.Yes),
			appconfig.CONFIG_UART_CONSOLE.Required(appconfig.Yes),
			appconfig.CONFIG_UART_LINE_CTRL.Required(appconfig.Yes),
			appconfig.CONFIG_PRINTK.Required(appconfig.Yes),

			appconfig.NewValue("CONFIG_ZBOSS_HALT_ON_ASSERT").Default(appconfig.Yes),
			appconfig.NewValue("CONFIG_RESET_ON_FATAL_ERROR").Default(appconfig.No),
			appconfig.NewValue("CONFIG_DEBUG_OPTIMIZATIONS").Default(appconfig.Yes),
			appconfig.NewValue("CONFIG_DEBUG_THREAD_INFO").Default(appconfig.Yes),
			appconfig.NewValue("CONFIG_THREAD_NAME").Default(appconfig.Yes),

			// Configurations for (hopefully) generating
			// good address for addr2line on exception.
			appconfig.NewValue("CONFIG_DEBUG_COREDUMP").Default(appconfig.Yes),
			appconfig.NewValue("CONFIG_DEBUG_COREDUMP_BACKEND_LOGGING").Default(appconfig.Yes),
			appconfig.NewValue("CONFIG_COREDUMP_DEVICE").Default(appconfig.Yes),
			// appconfig.NewValue("CONFIG_SYSTEM_WORKQUEUE_STACK_SIZE").Default(`2048`),
			// appconfig.NewValue("CONFIG_HEAP_MEM_POOL_SIZE").Default(`2048`),

			// ZBHome Debug enable
			appconfig.NewValue("CONFIG_ZBHOME_DEBUG_ENABLE").Required(appconfig.Yes),

			// Leds
			appconfig.NewValue("CONFIG_ZBHOME_DEBUG_LEDS").Required(ledsEnabled),

			// Log console
			appconfig.NewValue("CONFIG_ZBHOME_DEBUG_CONSOLE").Required(config.Console).Quoted(),
		},
		OverlayFn: overlayFn(config.Console),
	}
}

func (d *DebugConfig) IsEnabled() bool {
	return d != nil && d.Enabled
}

func overlayFn(console string) func(*devicetree.DeviceTree) error {
	if console == DebugConsoleUSB {
		console = "cdc_acm_uart0"
	}

	return func(dt *devicetree.DeviceTree) error {
		chosen := dt.FindSpecificNode(
			devicetree.SearchByName(devicetree.NodeNameRoot),
			devicetree.SearchByName(devicetree.NodeNameChosen))

		chosen.Properties = append(chosen.Properties,
			devicetree.NewProperty("zephyr,console", devicetree.Label(console)),
			devicetree.NewProperty("zephyr,shell-console", devicetree.Label(console)))

		return nil
	}
}
