---
title: Using the CLI
---

Zigbee_home CLI provides a way to do necessary actions to produce a device that is flashed with firmware generated & built with this project.

There are 3 things that CLI is responsible for:
* Generating the source code
* Building the firmware
* Flashing built firmware

Other functionalities might be available in the future.

## CLI commands

#### Global flags

At top level there is only option to provide path for configuration file, if it is not located in working directory.
```
--config <path_to_file>
```

Default configuration file called `zigbee.yml`.

For example

```bash
$ go run ./cmd/zigbee/... --config ../zigbee.yml
$ # Or if used with installed CLI
$ zigbee --config ../zigbee.yml
```

### `firmware`
This is a base command for sub-commands, and so it is not useful on it's own.
```bash
zigbee firmware
```

#### Flags
`--workdir`

:   Provide an optional working directory. This is where source code will be generated and used from.

    Defaults to current working directory.


#### Generate source code
Currently, generation of source code is not a separate action, as it is a part of building the firmware. A project goal is to do single call to do all the necessary actions, where user provides a configuration file and after running the CLI - a device will be flashed with generated firmware.

To generate source code from configuration file user must call 
```
firmware build --only-generate
```

For example
```bash
$ zigbee --config ../zigbee.yml firmware build --only-generate
```

#### Build the firmware
!!! warning
    Currently each time the build is requested - CLI will re-generate the source code. This means that any manual changes to the generated source code would be removed.

!!! note
    As building firmware requires having environment set correctly it is not possible to simply run this command from normal shell environment.
    If necessary this command can be executed from terminal that nRF Connect VS Code extension can create by going to `nRF Connect` extension -> in `Applications` right click necessary build configuration and select `Start new terminal here`.
    This will open a new terminal window with correct environment, in which this command can be executed.

To generate and build the firmware use
```
firmware build
```

This command does not need to be executed inside the generate source directory. Instead `--workdir` flag can be provided to point to necessary directory with source code.

For example
```bash
$ go run ./cmd/zigbee/... firmware --workdir ~/firmware/soil_moisture_sensor build
```

### Flashing the firmware
To flash the firmware user needs to set board to bootloader mode and then run
```
firmware flash
```
For example
```bash
$ zigbee firmware --workdir ~/firmware/soil_moisture_sensor flash
```

## Full example
This is a full example of what a typical flow for creating a configuration, generating & building source code and then flashing the device.

Example makes some assumptions:

- Board is nRF52840 Dongle with stock bootloader
- Generated source code will be located in `./firmware`

```console
$ cat zigbee.yml
general:
  runevery: 2m
  board: nrf52840dongle_nrf52840
  flasher: nrfutil
  flasheroptions:
    port: /dev/ttyACM1
board:
  i2c:
    - id: i2c0
      port: 0
      sda: 29
      scl: 31
sensors:
  - type: scd4x
    i2c:
      id: i2c0
$ ls
zigbee.yml
$ zigbee firmware --workdir ./firmware build --only-generate
$ # Now user goes to VS Code, opens the ./firmware directory and compiles the source code
$ # See 'Building the firmware' page for instructions on how to do it.
$ # As a verification that the firmware was built correctly execute following command:
$ ls ./firmware/build/zephyr/zephyr.hex
zephyr.hex
$ # User starts bootloader mode on the board by clicking reset button twice
$ zigbee firmware --workdir ./firmware flash
```