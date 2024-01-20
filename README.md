# Zigbee Home

Project that aims to provide similar functionality to [ESPHome](https://github.com/esphome/esphome), but for Zigbee devices.

# :information_source:  Note
`dev` branch is for experiments and exploration. 
It cannot be used to determine quality of resulting project.

## Status

Currently work is being carried to develop CLI application and adding sensors.

### Source Generation
CLI can generate source, which can then be built and flashed.

"Source" includes C source code, app config(`proj.conf`) and overlay (`app.overlay`).

### Building
CLI can't build the application yet, as building it requires environment to be properly set up (nrf- and zephyr-specific).

### Flashing
CLI can flash already built application with a couple of methods:
- nrfutil
- mcuboot
- west

### Supported devices

This project is being developed based on [nRF52840 Dongle](https://www.nordicsemi.com/Products/Development-hardware/nrf52840-dongle) from [Nordic Semiconductors](https://www.nordicsemi.com/)

Initial goal of the project is to support nRF52840 based devices, with expansion to nRF53 series.

## Using in Home Assistant
[ZHA](https://www.home-assistant.io/integrations/zha/) integration can figure out device capabilities without pre-defined configuration.

Firmware provided by this project is already working on the Dongle and reporting defined values when connected through ZHA.

[Zigbee2MQTT](https://www.zigbee2mqtt.io/) [recently](https://github.com/Koenkk/zigbee2mqtt/releases/tag/1.35.0) added ability to generate definitions for unknown devices, as ZHA does.
While it may not support all clusters and functionalities yet, there is active work going on to make support better.

## CLI

Go CLI application available in `/cli` will provide 
necessary options to build and upload firmware based on provided configuration.

For this to work user would need to already have nRF Connect SDK set up, with `west` and probably some flash tool applications available and working.

Defined flash tools are:
* `nrfutil` - already working
* `mcuboot`
* `west`

Default configuration file called `zigbee.yml` will be loaded 
and used to configure. 
See bare example in `cli/zigbee.yml`, or original configuration definition
in `cli/config/device.go:Device`.

To flash the board with built firmware run
`go run ./cli/cmd/zigbee firmware --workdir <path_to_project> flash`

Users can flash built applications with [nRF Connect for Desktop](https://www.nordicsemi.com/Products/Development-tools/nRF-Connect-for-Desktop) 
as well, if CLI tool is not suited for some cases.

This project needs nRF Connect SDK version 2.5.0. Other versions are not yet tested.

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