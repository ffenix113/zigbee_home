# This file is based on 
# https://github.com/martelmy/NCS_examples/blob/3cd874c68e13565ca89a082c67af54f48b704192/zigbee/light_bulb_dongle/pm_static_nrf52840dongle_nrf52840.yml
# , based on the answer here:
# https://devzone.nordicsemi.com/f/nordic-q-a/87169/zigbee-example-for-nrf52840-dongle/364086
#
# Thank you Marte Myrvold(https://github.com/martelmy)

{{- $flash := .AdditionalContext.Flash }}
{{- $sram := .AdditionalContext.SRAM }}

EMPTY_0:
  address: {{formatHex $flash.Address}}
  end_address: {{ formatHex $flash.EndAddress }}
  region: flash_primary
  size: {{ formatHex $flash.Size }}
EMPTY_1:
  address: {{ formatHex $sram.Address }}
  end_address: {{ formatHex $sram.EndAddress }}
  region: sram_primary
  size: {{ formatHex $sram.Size }}
