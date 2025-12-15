package appconfig

import "path/filepath"

// NOTE: This values are fetched from original `prj.conf` of this project,
// and some were updated to disable/remove some things as default.
// As such they do not represent good/best configurations,
// but mostly the ones that work for this project.
//
// Note: The names of values are not Pascal cased to signify that they are config options.
// The casing will change in future chore commits.
var (
	// ZBHome config
	CONFIG_ZBHOME_WATCHDOG_ENABLE = NewValue("CONFIG_ZBHOME_WATCHDOG_ENABLE").Default(Yes)

	// Zephyr config
	CONFIG_CPP                  = NewValue("CONFIG_CPP").Default(Yes)
	CONFIG_REQUIRES_FULL_LIBCPP = NewValue("CONFIG_REQUIRES_FULL_LIBCPP").Default(Yes)
	CONFIG_STD_CPP17            = NewValue("CONFIG_STD_CPP17").Default(Yes)

	// Logging
	CONFIG_CONSOLE          = NewValue("CONFIG_CONSOLE").Default(No)
	CONFIG_UART_CONSOLE     = NewValue("CONFIG_UART_CONSOLE").Default(No)
	CONFIG_UART_LINE_CTRL   = NewValue("CONFIG_UART_LINE_CTRL").Default(Yes)
	CONFIG_LOG_BACKEND_UART = NewValue("CONFIG_LOG_BACKEND_UART").Default(Yes)

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

	// For BLE OTA support.
	CONFIG_BT                           = NewValue("CONFIG_BT").Default(Yes)
	CONFIG_BT_PERIPHERAL                = NewValue("CONFIG_BT_PERIPHERAL").Default(Yes)
	CONFIG_BT_DEVICE_NAME               = NewValue("CONFIG_BT_DEVICE_NAME").Quoted().Default("ZBHome Device")
	CONFIG_NCS_SAMPLE_MCUMGR_BT_OTA_DFU = NewValue("CONFIG_NCS_SAMPLE_MCUMGR_BT_OTA_DFU").Default(Yes)

	// MCUBoot and firmware upgrade.
	CONFIG_BOOTLOADER_MCUBOOT             = NewValue("CONFIG_BOOTLOADER_MCUBOOT").Default(Yes)
	CONFIG_MCUBOOT_SIGNATURE_KEY_FILE     = NewValue("CONFIG_MCUBOOT_SIGNATURE_KEY_FILE").Quoted().Default(filepath.Join("\\${APPLICATION_CONFIG_DIR}", "mcuboot_key.pem"))
	CONFIG_MCUMGR_GRP_IMG_STATUS_HOOKS    = NewValue("CONFIG_MCUMGR_GRP_IMG_STATUS_HOOKS").Default(Yes)
	CONFIG_MCUMGR_GRP_OS_RESET_HOOK       = NewValue("CONFIG_MCUMGR_GRP_OS_RESET_HOOK").Default(Yes)
	CONFIG_MCUMGR_MGMT_NOTIFICATION_HOOKS = NewValue("CONFIG_MCUMGR_MGMT_NOTIFICATION_HOOKS").Default(Yes)

	// Settings
	CONFIG_ZMS          = NewValue("CONFIG_ZMS").Default(Yes)
	CONFIG_SETTINGS     = NewValue("CONFIG_SETTINGS").Default(Yes)
	CONFIG_SETTINGS_ZMS = NewValue("CONFIG_SETTINGS_ZMS").Default(Yes)

	// Sensors
	CONFIG_DHT = NewValue("CONFIG_DHT").Default(Yes)
)

// Sysbuild configuration options
var (
	// MCUBoot
	SB_CONFIG_BOOTLOADER_MCUBOOT          = NewValue("SB_CONFIG_BOOTLOADER_MCUBOOT").Default(Yes)
	SB_CONFIG_BOOT_SIGNATURE_TYPE_ED25519 = NewValue("SB_CONFIG_BOOT_SIGNATURE_TYPE_ED25519").Default(Yes)
	SB_CONFIG_BOOT_SIGNATURE_KEY_FILE     = NewValue("SB_CONFIG_BOOT_SIGNATURE_KEY_FILE").Quoted().Default(filepath.Join("\\${APPLICATION_CONFIG_DIR}", "mcuboot_key.pem"))
)

const (
	Yes = "y"
	No  = "n"
)
