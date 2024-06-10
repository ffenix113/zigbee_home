## Factory reset button

This example demonstrates how to set up custom factory reset button.

Note: if factory reset button is not defined in the configuration it will be the first button available in board definition. For example for nRF52840 Dongle it would be [this](https://github.com/zephyrproject-rtos/zephyr/blob/453ab8a9a356acf475a965a777a370795effa255/boards/nordic/nrf52840dongle/nrf52840dongle_nrf52840.dts#L64) button, for Adafruit Feather nRF52840 Express it would be [this](https://github.com/zephyrproject-rtos/zephyr/blob/453ab8a9a356acf475a965a777a370795effa255/boards/adafruit/feather/adafruit_feather_nrf52840.dts#L43) button. 

If board definition does not have any buttons and configuration does not define one - factory reset functionality through button press will not be available. In this case to reset Zigbee connection information it is needed to erase the device.

To properly test this example the board should be connected to a Zigbee network, and then the factory reset procedure(see `Usage`) should be done.

### Explanation
Factory reset will allow device to "forget" current network and start searching for new network.

For example when network was changed(encryption, channel, etc), but device still expects old network to be available. 
This will result in a situation where device will never connect to a network that it has stored in its memory.

### Usage

To actually trigger factory reset the user must hold the factory reset button for at least 5 seconds.
If debug & LEDs are enabled - LED for connection should then turn off and device automatically will start trying to connect to any open Zigbee networks.