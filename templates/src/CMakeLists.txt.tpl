#
# Copyright (c) 2022 Nordic Semiconductor ASA
#
# SPDX-License-Identifier: LicenseRef-Nordic-5-Clause
#

list(APPEND ZEPHYR_EXTRA_MODULES
  {{ range .Extenders }}{{ with .ZephyrModules }}{{range . -}}
  ${CMAKE_CURRENT_SOURCE_DIR}/modules/{{.}}
  {{- end }}{{ end }}{{ end }}
  )

cmake_minimum_required(VERSION 3.20.0)

################################################################################

# The application uses the configuration/<board> scheme for configuration files.
set(APPLICATION_CONFIG_DIR "${CMAKE_CURRENT_SOURCE_DIR}")

find_package(Zephyr REQUIRED HINTS $ENV{ZEPHYR_BASE})
project(zigbee_common)

################################################################################

# NORDIC SDK APP START
FILE(GLOB app_sources src/*.c src/**/*.c)
target_sources(app PRIVATE ${app_sources})

# NORDIC SDK APP END
