CLI relies heavily on the configuration file to do most actions.

For example the configuration file will provide information such as:

* Flashing configuration
* Whether the device should be router or not
* What peripherals are available on the board
* How frequently to poll the sensors
* What sensors are attached

By default, configuration file is called `zigbee.yml` and CLI will try to find it in the working directory. 
It is possible to point the CLI to this configuration file by providing a flag. For this see [Using the CLI](./index.md)

## Example configuration
!!! note
    The complete and latest file is located in the repository, in [`cli/zigbee.yml`](https://github.com/ffenix113/zigbee_home/blob/develop/cli/zigbee.yml). It may contain newer/changed configuration options.

```yaml
# Format for this file is not stable, and can change at any time.

# General configuraion that will be used by CLI
general:
  # Defines how much time a loop will sleep after each iteration.
  # Note: This does not necesseraly mean that each 2 minutes
  # sensor values will be updated. Currently all sensors are polled
  # synchronously, which means that if there are 3 sensors defined
  # and each needs 5 seconds to obtain the result then it will take
  # (3 * 5)s + runevery, which in this case will be 1m15s.
  runevery: 1m
  # Board value for now is used for building the source code
  # and possibly by flasher.
  board: nrf52840dongle_nrf52840
  # Flasher tells which flashing method to use.
  # Currently `nrfutil`, `mcuboot` and `west`
  # are defined(but not equally tested). Nrfutil works though.
  flasher: nrfutil
  # Flasheroptions are flasher-specific options.
  flasheroptions:
    port: /dev/ttyACM1

# This section is for defining peripherals
# on the board. I.e. uart, spi, i2c, etc.
# NOTE: Only changes should be defined here.
# See https://github.com/zephyrproject-rtos/zephyr/tree/main/boards/<arch>/<board_name>/<board_name>.dts
# for existing definitions for the board.
# For example nRF52840 Dongle would have board devicetree at
# https://github.com/zephyrproject-rtos/zephyr/tree/main/boards/arm/nrf52840dongle_nrf52840/nrf52840dongle_nrf52840.dts
board:
  # This option will add USB UART loging functionality.
  # Quite limited for now, but can be
  # extended through template modifications.
  debuglog: false
  # Change device type from sleepy end device to router.
  # This allows for device to always be available from coordinator
  # and other devices.
  # By default device will be configured as sleepy end device.
  is_router: true
  # I2C is optional, but strongly suggested if some sensors use it. 
  # It is used to define enabled I2C buses
  # and optionally different pins.
  i2c:
    # ID of instance is the same as defined in the SoC definition.
    # Generally they are in form of `i2c[0-9]`.
    # Number of i2c instances is also limited, and this allows only to
    # re-define pins for specified I2C instance.
    - id: i2c0
      port: 0
      sda: 29
      scl: 31
    # Pins can be omitted if a specific I2C instance should only be enabled.
    # Then it will use default pins for this bus.
    - id: i2c1

# Sensors define a list of devices that 
# can provide sensor values or be controlled
sensors:
  # All sensors have type, and most will 
  # also have sensor-specific configuration.
  - type: bme680
    i2c:
      id: i2c0
      # Some devices might have changable I2C address, 
      # which can be defined here.
      # Note: this does not change the device address,
      # only tells which address to use.
      addr: '0x76'
  - type: scd4x
    i2c:
      id: i2c0
  # - type: device_temperature
  # on_off is a sensor that will respond to on/off state of the pin.
  # For now verifyied to be controlled by the client only,
  # so not by actually changing the state of the pin manually.
  # - type: on_off
  #   pin:
  #     # This is Green LED(LD1) on nrf52840 Dongle
  #     port: 0
  #     pin: 6
  #     inverted: true
```

## Sensor options
Each sensor can have unique configuration options. To define such options they must be provided in their configuration in `sensors` section of the configuration file.

If sensors do not provide required configuration - the generated source code might not compile with some errors, or sensor might not report values correctly.

Some sensors would not require any configuration:
```yaml
sensors:
  - type: device_temperature
```
, while others might need some:
```yaml
sensors:
  - type: bme680
    i2c:
      id: i2c0
      addr: '0x76'
```
Here we define that the board will have a Bosch BME680 sensor attached to I2C instance with id `i2c0`, and will have address `0x76`.

Other options may also be specified in the future, for example:
```yaml
sensors:
  - type: bme680
    temperature:
      oversampling: x4 # Notice the `oversampling` option
    i2c:
      id: i2c0
      addr: '0x76'
```