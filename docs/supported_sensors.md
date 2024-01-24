## Supporting more sensors
Depending on a sensor it might be either easy, or rather involving to provide support for a sensor.

Currently there is a work to change how sensors are polled to make all sensors to work through Zephyr's Sensor API, which would result in same codebase for all sensors in the main app, while drivers will implement interface required to talk to the sensors.

When this change will be done - all sensors officially supported by Zephyr, will be supported. Note that the sensor will be supported only in the context of [Zigbee Cluster Library](https://zigbeealliance.org/wp-content/uploads/2021/10/07-5123-08-Zigbee-Cluster-Library.pdf), meaning that if sensor supports some measurements that Zigbee Cluster Library does not support - it will not be available in Zigbee.

## Supported sensors
* Bosch BME280 / BME680
* Sensirion SCD41 ([driver](https://github.com/nobodyguy/sensirion_zephyr_drivers) by [nobodyguy](https://github.com/nobodyguy))