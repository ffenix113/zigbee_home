In some situations it is necessary to have a way to disconnect device from the network it was connected to.

For example:
* Device will be used in different network.
* Previous network is not available (i.e. on configuration change).
* Device was not correctly disconnected from network, resulting in device having network configuration, but network does not accept the device.

To resolve this issues firmware provides a way to set a button that will act as a factory reset button:
```yml
board:
    factory_reset_button: btn1
    buttons:
        - id: btn1
```

To perform a factory reset user needs to hold `btn1` for more than 5 seconds and then release.
After this the device will remove current network configuration and automatically try to pair with any open Zigbee network.