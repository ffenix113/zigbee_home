# SPDX-License-Identifier: Apache-2.0

if(CONFIG_SCD4X)
    zephyr_include_directories(.)

    zephyr_library()
    zephyr_library_sources(scd4x.c)
endif()