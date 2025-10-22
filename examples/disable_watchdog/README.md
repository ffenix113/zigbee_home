## Disable watchdog

This example shows how to disable watchdog.

Generally users should not need to do this, but if some specific requirements need to disable it - it can be done.
For example in case of minimal possible power consumption.

Watchdog is always disabled when debug configuration is enabled.

### Notes
* Watchdog power consumption may be negligable compared to selected sensors and other options.
* If device locks-up or otherwise misbehaves and watchdog is disabled it would require manual reset / power cycle of the board.