## IAS Zone as contact sensor

This example provides definition for contact sensor. 

Each time the state of the pin `0.04` will change - board will update coordinator to show the new value.


## Note: analog pins only!

Please note that this sensor can work only with analog pins. See [pin assingment](https://infocenter.nordicsemi.com/topic/ps_nrf52840/pin.html?cp=5_0_0_6_0) and look for `AIN` pins, also see table below for same information:
* 0.02 - AIN0
* 0.03 - AIN1
* 0.04 - AIN2
* 0.05 - AIN3
* 0.28 - AIN4
* 0.29 - AIN5
* 0.30 - AIN6
* 0.31 - AIN7