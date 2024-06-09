## Factory reset button

This example demonstrates how to set up custom factory reset button.

To properly test this example the board should be connected to a Zigbee network, and then the factory reset procedure(see `Usage`) should be done.

### Explanation
Factory reset will allow device to "forget" current network and start searching for new network.

For example when network was changed(encryption, channel, etc), but device still expects old network to be available. 
This will result in a situation where device will never connect to a network that it has stored in its memory.

### Usage

To actually trigger factory reset the user must hold the factory reset button for at least 5 seconds.
If debug & LEDs are enabled - LED for connection should then turn off and device automatically will start trying to connect to any open Zigbee networks.