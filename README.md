# Zigbee Home

Project that aims to provide similar functionality to [ESPHome](https://github.com/esphome/esphome), but for Zigbee devices.

# :information_source:  Note
`dev` branch is for experiments and exploration. 
It cannot be used to determine quality of resulting project.

## Status

Currently work is being carried to develop CLI application and adding sensors.

Priorities can be ordered as:
* Board(bootloader) support
* Adding known sensors
* Adding Zigbee clusters & templates for unavailable clusters in ZBOSS

### Licenses
This project uses information from ZBOSS SDK, license for which can be found in `zboss_license.txt`. 

### References
* nRF Connect SDK
* * [Download page](https://www.nordicsemi.com/Products/Development-software/nRF-Connect-SDK)
* * [Source](https://github.com/nrfconnect/sdk-nrf)
* * [Documentation](http://developer.nordicsemi.com/nRF_Connect_SDK/doc/latest)
* [Zephyr project](https://www.zephyrproject.org/)
* [ESPHome](https://esphome.io/)
* [Zigbee Cluster Library](https://csa-iot.org/wp-content/uploads/2022/01/07-5123-08-Zigbee-Cluster-Library-1.pdf)

## Special thanks
* @rsporsche - for donating nRF52840 DK board
* @Hedda - for informational support
