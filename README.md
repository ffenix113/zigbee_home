# Zigbee Home

Project that aims to provide similar functionality to [ESPHome](https://github.com/esphome/esphome), but for Zigbee devices.

# :information_source:  Note
`dev` branch is for experiments and exploration. 
It cannot be used to determine quality of resulting project.

## Status

Extremely early development, no working parts yet.

Currently work is being carried to develop CLI application.

### Supported devices

This project is being developed based on [nRF52840 Dongle](https://www.nordicsemi.com/Products/Development-hardware/nrf52840-dongle) from [Nordic Semiconductors](https://www.nordicsemi.com/)

Initial goal of the project is to support nRF52840 based devices, with expansion to nRF53 series.

## Using in Home Assistant
[ZHA](https://www.home-assistant.io/integrations/zha/) integration can figure out device capabilities without pre-defined configuration.

Firmware provided by this project is already working on the Dongle and reporting defined values when connected through ZHA.

[Zigbee2MQTT](https://www.zigbee2mqtt.io/) requires to define device configuration before it can be useful, so for quick set-up it is not as useful.

## CLI

Go CLI application available in `/cli` will provide necessary options to build and upload firmware based on provided configuration.

For this to work user would need to already have nRF Connect SDK set up, with `west` and either [nrfutil](https://www.nordicsemi.com/Products/Development-tools/nRF-Util) or [nRF Connect for Desktop](https://www.nordicsemi.com/Products/Development-tools/nRF-Connect-for-Desktop) applications available and working.

This project needs nRF Connect SDK version 2.5.0. Other versions are not yet tested.

### References
* nRF Connect SDK
* * [Download page](https://www.nordicsemi.com/Products/Development-software/nRF-Connect-SDK)
* * [Source](https://github.com/nrfconnect/sdk-nrf)
* * [Documentation](http://developer.nordicsemi.com/nRF_Connect_SDK/doc/latest)
* [Zephyr project](https://www.zephyrproject.org/)
* [ESPHome](https://esphome.io/)