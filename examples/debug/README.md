## Debug configuration

Enables simple, more debug-friendly configuration, that will provide more information to the user on the behavior.

This includes
- LEDs that will show power & connection state
- USB interface that will print logs.

Configuration defines `led1` as "forwarded" from the board configuration. `led2` is defined as LED connected to pin `0.04`.

### Reading USB logs

On Linux(tested on Ubuntu) the logs can be read using programs like `screen` and `minicom`.

For example with `minicom` the logs can be read with: `sudo minicom -D /dev/ttyACM0 115200`. `ACM0` can be different on your machine depending on the configuration and connected devices.
Available `ACM` interfaces can be checked with `ls /dev/ttyA*`. When board is connected - one of the `ACM` devices will be boards logging output.