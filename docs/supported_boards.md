## Plans for supporting boards
The plan for board support is to have as much nrf52840-based boards to make compatible as possible.

This probably can be achieved by utilizing MCUboot, to make the bootloader the same on all boards.
It will also bring possible benefits like serial flashing, DFU through BLE and OTA for Zigbee.

### What is the problem? Why MCUboot?

As of now different boards can come with a couple of bootloaders, each one of them may require different approach for flashing and different set of configuration options. Supporting each possible configuration may not be easy, nor needed.

With MCUboot the interface between boards and zigbee_home should be the same, independednt of the board.

### MCUboot support
There is on-going work to investigate configuration and usage of MCUboot as second-stage bootloader for nRF52840 Dongle(as it is shipped with nRF5 Bootloader), which should lead a way to support on other boards. For example support for another board is work-in-progress based on MCUboot: [Arduino nano 33 BLE](https://store.arduino.cc/products/arduino-nano-33-ble). 

### If not MCUboot as second-stage bootloader?
If MCUboot would not be enough to make boards work correctly and support necessary functionality - then another solution could be to flash it directly via JTag/SWD, in hopes that this will make it work.

## Currently supported boards
This list can be better thought of as "supported bootloaders" instead of the boards.

* nRF52840 Dongle (`nrf52840dongle_nrf52840`) - Uses `nrf52_legacy` bootloader.
* Arduino Nano 33 BLE (Sense) (`arduino_nano_33_ble`) - Uses `arduino` bootloader. With this configuration other boards with nrf52840 and Bossac bootloader(== Arduino bootloader) might be supported. Give it a try!

### Experimental
* Any Adafruit bootloader-based boards.
:   This includes support for bootloaders with SD S132, S140 v6 and v7. A `bootloader` configuration option must be set to a necessary bootloader. Supported bootloader versions are
    * `adafruit_nrf52_sd132`
    * `adafruit_nrf52_sd140_v6`
    * `adafruit_nrf52_sd140_v7`