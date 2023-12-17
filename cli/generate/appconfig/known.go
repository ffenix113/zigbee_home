package appconfig

// NOTE: This values are fetched from original `prj.conf` of this project.
// As such they do not represent good/best configurations,
// but mostly the ones that work for this project.
var (
	// Logging
	CONFIG_LOG              = NewValue("CONFIG_LOG").Default(Yes)
	CONFIG_SERIAL           = NewValue("CONFIG_SERIAL").Default(Yes)
	CONFIG_CONSOLE          = NewValue("CONFIG_CONSOLE").Default(Yes)
	CONFIG_UART_CONSOLE     = NewValue("CONFIG_UART_CONSOLE").Default(Yes)
	CONFIG_UART_LINE_CTRL   = NewValue("CONFIG_UART_LINE_CTRL").Default(Yes)
	CONFIG_LOG_BACKEND_UART = NewValue("CONFIG_LOG_BACKEND_UART").Default(No)

	// USB
	CONFIG_USB_DEVICE_INITIALIZE_AT_BOOT = NewValue("CONFIG_USB_DEVICE_INITIALIZE_AT_BOOT").Default(No)
	CONFIG_USB_DEVICE_PRODUCT            = NewValue("CONFIG_USB_DEVICE_PRODUCT").Default(`"Dongle: Zigbee Device"`)
	CONFIG_USB_DEVICE_PID                = NewValue("CONFIG_USB_DEVICE_PID").Default(`0x0004`)
	CONFIG_USB_DEVICE_STACK              = NewValue("CONFIG_USB_DEVICE_STACK").Default(Yes)

	// Drivers / peripherals
	CONFIG_I2C        = NewValue("CONFIG_I2C").Default(Yes)
	CONFIG_SENSOR     = NewValue("CONFIG_SENSOR").Default(Yes)
	CONFIG_BME280     = NewValue("CONFIG_BME280").Default(Yes).Depends(CONFIG_I2C, CONFIG_SENSOR)
	CONFIG_DK_LIBRARY = NewValue("CONFIG_DK_LIBRARY").Default(Yes)

	// Zigbee
	CONFIG_ZIGBEE                 = NewValue("CONFIG_ZIGBEE").Default(Yes)
	CONFIG_ZIGBEE_APP_UTILS       = NewValue("CONFIG_ZIGBEE_APP_UTILS").Default(Yes)
	CONFIG_ZIGBEE_CHANNEL         = NewValue("CONFIG_ZIGBEE_CHANNEL").Default(`11`)
	CONFIG_ZIGBEE_ROLE_END_DEVICE = NewValue("CONFIG_ZIGBEE_ROLE_END_DEVICE").Default(Yes)

	// Cryptography
	CONFIG_CRYPTO               = NewValue("CONFIG_CRYPTO").Default(Yes)
	CONFIG_CRYPTO_NRF_ECB       = NewValue("CONFIG_CRYPTO_NRF_ECB").Default(Yes).Depends(CONFIG_CRYPTO)
	CONFIG_CRYPTO_INIT_PRIORITY = NewValue("CONFIG_CRYPTO_INIT_PRIORITY").Default(`80`)

	// Power configuration
	CONFIG_RAM_POWER_DOWN_LIBRARY = NewValue("CONFIG_RAM_POWER_DOWN_LIBRARY").Default(Yes)

	// Network
	CONFIG_NET_IPV6          = NewValue("CONFIG_NET_IPV6").Default(No)
	CONFIG_NET_IP_ADDR_CHECK = NewValue("CONFIG_NET_IP_ADDR_CHECK").Default(No)
	CONFIG_NET_UDP           = NewValue("CONFIG_NET_UDP").Default(No)

	// Debug
	CONFIG_ZBOSS_HALT_ON_ASSERT        = NewValue("CONFIG_ZBOSS_HALT_ON_ASSERT").Default(Yes)
	CONFIG_RESET_ON_FATAL_ERROR        = NewValue("CONFIG_RESET_ON_FATAL_ERROR").Default(No)
	CONFIG_SYSTEM_WORKQUEUE_STACK_SIZE = NewValue("CONFIG_SYSTEM_WORKQUEUE_STACK_SIZE").Default(`2048`)
	CONFIG_HEAP_MEM_POOL_SIZE          = NewValue("CONFIG_HEAP_MEM_POOL_SIZE").Default(`2048`)
	CONFIG_DEBUG_OPTIMIZATIONS         = NewValue("CONFIG_DEBUG_OPTIMIZATIONS").Default(No)
	CONFIG_DEBUG_THREAD_INFO           = NewValue("CONFIG_DEBUG_THREAD_INFO").Default(No)
)

const (
	Yes = "y"
	No  = "n"
)
