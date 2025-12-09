## BLE OTA

This example shows how to set up BLE OTA so that board can be flashed through BLE.

BLE OTA will allow users to update device firmware wirelessly.
It will be possible to remove BLE OTA if necessary by removing necessary configuration
and updating the firmware.

It uses [SMP](https://docs.zephyrproject.org/3.7.0/services/device_mgmt/smp_protocol.html) protocol to do so, 
so it is compatible with any client applications that support this protocol.

### How to update firmware manually (without `zigbee_home`)
If, for some reason, it is not possible to use `zigbee_home` to flash the firmware - 
the update process is designed to not rely on `zigbee_home` or "proprietary" code.

The steps are:
1. Build the firmware (using correct firmware signing key)
2. Request board reset with `force` option set to true via SMP.
    This step is optional, but it will significantly increase the speed of firmware upload.
    When firmware will receive this request - it will disable Zigbee stack on reboot, 
    so only BLE will be using device radio.
3. Upload file `zephyr.signed.bin` via SMP [image upload](https://docs.zephyrproject.org/3.7.0/services/device_mgmt/smp_groups/smp_group_1.html#image-upload).
4. After full image transfer device will mark new image to be tested on reboot
5. Request board reset via SMP
6. MCUBoot will now verify new image signature and move it to the first slot.
    This step will take ~2-3 minutes currently, and it is expected.
7. After image was successfully moved - it will be loaded and marked as permanent.

Steps 2, 3 and 5 are the ones that use SMP protocol.

`zigbee_home` actually executes exactly these steps to update the firwmare, so user does not have to.

### Notes
- Board must be flashed via SWD/JLink the first time the project will be built, as it will contain specific cryptographic key.
    And also MCUBoot must be flashed as bootloader (with specific cryptographic key).
- Energy usage will increase, as BLE will always be on. Even though BLE has reduced advertisement interval.
- Currently anyone will be able to see and communicate with BLE device. Communication is not authenticated, nor encrypted.
    Both features will be coming in the future, as long as it will be possible to implement them both
    for firmware- and client-side. So board must be able to authenticate and/or encrypt connection, 
    but client (`zigbee_home` application, for example) must also be able to do the same.
    
    If `experimental.bleota.key` is set in configuration - it will not be possible to trick device
    into loading incorrect firmware. As such this option is highly recommended, 
    and the key value must be unique per device.
- If the OTA key changes - new firmware must be flashed via SWD/JLink.
- Size of the firwmare will be increased, so it may be possible that currently some devices will not be able to flash it.
    In this case please add discussion with your configuration and board.
- It is not required to have `zigbee_home` to flash the firmware. It only simplifies this work by doing necessary steps.