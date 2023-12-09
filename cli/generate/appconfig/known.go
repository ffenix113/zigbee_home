package appconfig

// NOTE: This values are fetched from original `prj.conf` of this project.
// As such they do not represent good/best configurations,
// but mostly the ones that work for this project.
var (
	// Logging
	CONFIG_LOG              = NewValue("CONFIG_LOG").Default(`y`)
	CONFIG_SERIAL           = NewValue("CONFIG_SERIAL").Default(`y`)
	CONFIG_CONSOLE          = NewValue("CONFIG_CONSOLE").Default(`y`)
	CONFIG_UART_CONSOLE     = NewValue("CONFIG_UART_CONSOLE").Default(`y`)
	CONFIG_UART_LINE_CTRL   = NewValue("CONFIG_UART_LINE_CTRL").Default(`y`)
	CONFIG_LOG_BACKEND_UART = NewValue("CONFIG_LOG_BACKEND_UART").Default(`n`)

	// USB
	CONFIG_USB_DEVICE_INITIALIZE_AT_BOOT = NewValue("CONFIG_USB_DEVICE_INITIALIZE_AT_BOOT").Default(`n`)
	CONFIG_USB_DEVICE_PRODUCT            = NewValue("CONFIG_USB_DEVICE_PRODUCT").Default(`"Dongle: Zigbee Weather Station"`)
	CONFIG_USB_DEVICE_PID                = NewValue("CONFIG_USB_DEVICE_PID").Default(`0x0004`)
	CONFIG_USB_DEVICE_STACK              = NewValue("CONFIG_USB_DEVICE_STACK").Default(`y`)

	// Drivers / peripherals
	CONFIG_I2C        = NewValue("CONFIG_I2C").Default(`y`)
	CONFIG_SENSOR     = NewValue("CONFIG_SENSOR").Default(`y`)
	CONFIG_BME280     = NewValue("CONFIG_BME280").Default(`y`)
	CONFIG_DK_LIBRARY = NewValue("CONFIG_DK_LIBRARY").Default(`y`)

	// Zigbee
	CONFIG_ZIGBEE                 = NewValue("CONFIG_ZIGBEE").Default(`y`)
	CONFIG_ZIGBEE_APP_UTILS       = NewValue("CONFIG_ZIGBEE_APP_UTILS").Default(`y`)
	CONFIG_ZIGBEE_CHANNEL         = NewValue("CONFIG_ZIGBEE_CHANNEL").Default(`11`)
	CONFIG_ZIGBEE_ROLE_END_DEVICE = NewValue("CONFIG_ZIGBEE_ROLE_END_DEVICE").Default(`y`)

	// Cryptography
	CONFIG_CRYPTO               = NewValue("CONFIG_CRYPTO").Default(`y`)
	CONFIG_CRYPTO_NRF_ECB       = NewValue("CONFIG_CRYPTO_NRF_ECB").Default(`y`)
	CONFIG_CRYPTO_INIT_PRIORITY = NewValue("CONFIG_CRYPTO_INIT_PRIORITY").Default(`80`)

	// Power configuration
	CONFIG_RAM_POWER_DOWN_LIBRARY = NewValue("CONFIG_RAM_POWER_DOWN_LIBRARY").Default(`y`)

	// Network
	CONFIG_NET_IPV6          = NewValue("CONFIG_NET_IPV6").Default(`n`)
	CONFIG_NET_IP_ADDR_CHECK = NewValue("CONFIG_NET_IP_ADDR_CHECK").Default(`n`)
	CONFIG_NET_UDP           = NewValue("CONFIG_NET_UDP").Default(`n`)

	// Debug
	CONFIG_ZBOSS_HALT_ON_ASSERT        = NewValue("CONFIG_ZBOSS_HALT_ON_ASSERT").Default(`y`)
	CONFIG_RESET_ON_FATAL_ERROR        = NewValue("CONFIG_RESET_ON_FATAL_ERROR").Default(`n`)
	CONFIG_SYSTEM_WORKQUEUE_STACK_SIZE = NewValue("CONFIG_SYSTEM_WORKQUEUE_STACK_SIZE").Default(`2048`)
	CONFIG_HEAP_MEM_POOL_SIZE          = NewValue("CONFIG_HEAP_MEM_POOL_SIZE").Default(`2048`)
	CONFIG_DEBUG_OPTIMIZATIONS         = NewValue("CONFIG_DEBUG_OPTIMIZATIONS").Default(`n`)
	CONFIG_DEBUG_THREAD_INFO           = NewValue("CONFIG_DEBUG_THREAD_INFO").Default(`n`)
)
