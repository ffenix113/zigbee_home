config ZBHOME_DEBUG_ENABLE
	bool "Debug configuration"
	default n
	help
		Enable debug configuration for the firmware

config ZBHOME_DEBUG_LEDS
	bool "Enable debug LEDs"
	default y
	depends on ZBHOME_DEBUG_ENABLE
	help
	  Enable LEDs to signify some debug state(on/off, joined Zigbee network).

config ZBHOME_DEBUG_CONSOLE
	string "Output console for logs"
	default ""
	help
		Specifies which backend to use for logging

source "Kconfig.zephyr"

module = ZIGBEE_DEVICE
module-str = Zigbee device