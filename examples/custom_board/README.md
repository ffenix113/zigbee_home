## Custom board support

This example shows how to make use of customboard definitions that are not (yet) supported by nRF Connect (and Zephyr) version that is being used.

For this example to work you must set working directory to be the one which contains `boards` directory with your custom board(s).

Note: If you will run `zigbee_home firmware build --workdir ${DIR} --clear-work-dir` - it will remove the `bords` directory as well! Please only do it if required and you know what the result will be.

### Creating custom board

The process is rather technical, so it may not be suited for everybody.

Zephyr documentation on board support can be found here: (https://docs.zephyrproject.org/latest/hardware/porting/board_porting.html)[https://docs.zephyrproject.org/latest/hardware/porting/board_porting.html]
