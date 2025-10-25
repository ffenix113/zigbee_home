config ZBHOME_WATCHDOG_ENABLE
	bool "Enable watchdog"
	default y
	imply WATCHDOG
	help
	  Enable watchdog that will reset SoC (and so the board) in case of a lock up.

config ZBHOME_DEBUG_ENABLE
	bool "Debug configuration"
	default n
	imply DEBUG
	imply DEBUG_THREAD_INFO
	imply THREAD_NAME
	imply EXCEPTION_STACK_TRACE

	# Logging
	imply LOG
	# Selecting instead of implying because it is
	# set to "n" in config.
	select CONSOLE
	imply SERIAL
	imply LOG_BACKEND_UART
	imply UART_CONSOLE
	imply UART_LINE_CTRL

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

config ZBHOME_WITH_MCUBOOT
	bool "Add MCUboot as bootloader"
	default n

choice ZBOSS_ASSERT_HANDLER
	default ZBOSS_HALT_ON_ASSERT
endchoice

# Configurations below should provide required configuraiton options for specific soc series.

config ZBHOME_NRF53X
	def_bool y
	imply NRF_SECURITY
	imply ZIGBEE_USE_SOFTWARE_AES
	depends on SOC_SERIES_NRF53X

config ZBHOME_NRF52X
	def_bool y
	imply CRYPTO
	imply CRYPTO_NRF_ECB
	depends on SOC_SERIES_NRF52X

# Always on configuration to set up always required configs.
config ZBHOME
	def_bool y
	imply ZIGBEE_ADD_ON

if !ZIGBEE_ADD_ON
config ZIGBEE
	default y
endif

source "Kconfig.zephyr"
source "${ZEPHYR_BASE}/../nrf/subsys/zigbee/Kconfig"

module = ZIGBEE_DEVICE
module-str = Zigbee device