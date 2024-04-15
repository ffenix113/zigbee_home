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