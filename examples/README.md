## Examples

This folder shows some examples on the ways to configure the device, bootloader and sensors.

Most of the examples can be combined, if they are not overwriting each-others configuration, or if user merges them with most appropriate options.

The examples may not provide configurations that will definitely run on your hardware configuration, instead providing examples of options that can be used to achieve something.

## Provided examples
* `basic` - the most featureless firmware that can be build with this project.
* `adafruit_bootloader` - shows how to specify different supported bootloader(in this case Adafruit) for the board.
* `sensor_bme280` - usage of Bosch BME280 sensoron I2C interface.
* `config_tags` - configuration tags that can help define more extensible/generic configuration.
* `debug` - simple debug configuration that will show power & connection state + USB logging.
* `contact` - configuration that uses IAS Zone to add a contact sensor.