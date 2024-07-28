## Set toolchain version

User can provide specific toolchain version that should be used.

By default, if it is not defined explicitely, the default version from
toolchain configuration will be used.

If version specified in the configuration file is not available - 
next available patch version will be selected.

Zigbee_home has default version of toolchain specified, 
that will be selected if no other version is specified 
in configuration file or environment. 
This default version is defined in `config/device.go`. 