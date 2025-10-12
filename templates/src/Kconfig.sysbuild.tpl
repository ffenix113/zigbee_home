source "${ZEPHYR_BASE}/share/sysbuild/Kconfig"

# Set 802.15.4 RPC for network core on nrf53 series.

config NRF_DEFAULT_802154 
    default y
    depends on SOC_SERIES_NRF53X