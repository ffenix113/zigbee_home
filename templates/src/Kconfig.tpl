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

config ZBHOME_ZIGBEE_CRYPTO
	bool
	default y
	# For nRF5340
	imply NRF_SECURITY if SOC_SERIES_NRF53X
	imply ZIGBEE_USE_SOFTWARE_AES if SOC_SERIES_NRF53X
	# For nRF52840
	imply CONFIG_CRYPTO if SOC_SERIES_NRF52X
	imply CONFIG_CRYPTO_NRF_ECB if SOC_SERIES_NRF52X
	help
	  "Set correct crypto configuration for specific chip"

source "Kconfig.zephyr"

module = ZIGBEE_DEVICE
module-str = Zigbee device