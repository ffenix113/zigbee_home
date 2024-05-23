package appconfig

// NOTE: This values are fetched from original `prj.conf` of this project,
// and some were updated to disable/remove some things as default.
// As such they do not represent good/best configurations,
// but mostly the ones that work for this project.
var (
	// Logging
	CONFIG_LOG              = NewValue("CONFIG_LOG").Default(No)
	CONFIG_SERIAL           = NewValue("CONFIG_SERIAL").Default(No)
	CONFIG_CONSOLE          = NewValue("CONFIG_CONSOLE").Default(No)
	CONFIG_UART_CONSOLE     = NewValue("CONFIG_UART_CONSOLE").Default(No)
	CONFIG_UART_LINE_CTRL   = NewValue("CONFIG_UART_LINE_CTRL").Default(Yes)
	CONFIG_LOG_BACKEND_UART = NewValue("CONFIG_LOG_BACKEND_UART").Default(Yes)
	CONFIG_PRINTK           = NewValue("CONFIG_PRINTK").Default(Yes)

	// USB
	CONFIG_USB_DEVICE_INITIALIZE_AT_BOOT = NewValue("CONFIG_USB_DEVICE_INITIALIZE_AT_BOOT").Default(No)
	CONFIG_USB_DEVICE_PRODUCT            = NewValue("CONFIG_USB_DEVICE_PRODUCT").Default(`"Dongle: Zigbee Device"`)
	CONFIG_USB_DEVICE_PID                = NewValue("CONFIG_USB_DEVICE_PID").Default(`0x0004`)
	CONFIG_USB_DEVICE_STACK              = NewValue("CONFIG_USB_DEVICE_STACK").Default(No)

	// Drivers / peripherals
	CONFIG_I2C        = NewValue("CONFIG_I2C").Default(Yes)
	CONFIG_SENSOR     = NewValue("CONFIG_SENSOR").Default(No)
	CONFIG_BME280     = NewValue("CONFIG_BME280").Default(No).Depends(CONFIG_I2C.Required(Yes), CONFIG_SENSOR.Required(Yes))
	CONFIG_DK_LIBRARY = NewValue("CONFIG_DK_LIBRARY").Default(Yes)

	// Zigbee
	CONFIG_ZIGBEE                              = NewValue("CONFIG_ZIGBEE").Default(Yes)
	CONFIG_ZIGBEE_APP_UTILS                    = NewValue("CONFIG_ZIGBEE_APP_UTILS").Default(Yes)
	CONFIG_ZIGBEE_CHANNEL_MASK                 = NewValue("CONFIG_ZIGBEE_CHANNEL_MASK").Default("0x7FFF800")
	CONFIG_ZIGBEE_ROLE_END_DEVICE              = NewValue("CONFIG_ZIGBEE_ROLE_END_DEVICE").Default(Yes)
	CONFIG_ZIGBEE_ROLE_ROUTER                  = NewValue("CONFIG_ZIGBEE_ROLE_ROUTER").Default(Yes)
	CONFIG_ZIGBEE_CHANNEL_SELECTION_MODE_MULTI = NewValue("CONFIG_ZIGBEE_CHANNEL_SELECTION_MODE_MULTI").Default(Yes)

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

	// Sensors
	CONFIG_DHT = NewValue("CONFIG_DHT").Default(Yes)
)

const (
	Yes = "y"
	No  = "n"
)
