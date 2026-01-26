/*
 * ZBOSS Zigbee 3.0
 *
 * Copyright (c) 2012-2020 DSR Corporation, Denver CO, USA.
 * www.dsr-zboss.com
 * www.dsr-corporation.com
 * All rights reserved.
 *
 *
 * Use in source and binary forms, redistribution in binary form only, with
 * or without modification, are permitted provided that the following conditions
 * are met:
 *
 * 1. Redistributions in binary form, except as embedded into a Nordic
 *    Semiconductor ASA integrated circuit in a product or a software update for
 *    such product, must reproduce the above copyright notice, this list of
 *    conditions and the following disclaimer in the documentation and/or other
 *    materials provided with the distribution.
 *
 * 2. Neither the name of Nordic Semiconductor ASA nor the names of its
 *    contributors may be used to endorse or promote products derived from this
 *    software without specific prior written permission.
 *
 * 3. This software, with or without modification, must only be used with a Nordic
 *    Semiconductor ASA integrated circuit.
 *
 * 4. Any software provided in binary form under this license must not be reverse
 *    engineered, decompiled, modified and/or disassembled.
 *
 * THIS SOFTWARE IS PROVIDED BY NORDIC SEMICONDUCTOR ASA "AS IS" AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY, NONINFRINGEMENT, AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL NORDIC SEMICONDUCTOR ASA OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR
 * TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
 * THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
/* PURPOSE:
 */

#pragma once

#ifdef __cplusplus
extern "C"
{
#endif

// #include <zboss_api_buf.h>
#include <zboss_api.h>
// #include <zb_zcl_el_measurement.h>
#include <zcl/zb_zcl_el_measurement.h>

#ifdef __cplusplus
}
#endif

// These values are based on the INA226 with 0.2 mOhm resistor
// to measure up to 20 A.
// FIXME: This values should not be static, but per-cluster instance.
#define ZCL_DC_VOLTAGE_VALUE_DIVISOR 500
#define ZCL_DC_CURRENT_VALUE_DIVISOR 1000
#define ZCL_DC_POWER_VALUE_DIVISOR 30

#define ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DC_VOLTAGE_ID(data_ptr) \
    {                                                                                     \
        ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DC_VOLTAGE_ID,                                 \
        ZB_ZCL_ATTR_TYPE_S16,                                                             \
        ZB_ZCL_ATTR_ACCESS_READ_ONLY | ZB_ZCL_ATTR_ACCESS_REPORTING,                      \
        (ZB_ZCL_NON_MANUFACTURER_SPECIFIC),                                               \
        (void *)data_ptr}

#define ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DC_CURRENT_ID(data_ptr) \
    {                                                                                     \
        ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DC_CURRENT_ID,                                 \
        ZB_ZCL_ATTR_TYPE_S16,                                                             \
        ZB_ZCL_ATTR_ACCESS_READ_ONLY | ZB_ZCL_ATTR_ACCESS_REPORTING,                      \
        (ZB_ZCL_NON_MANUFACTURER_SPECIFIC),                                               \
        (void *)data_ptr}

#define ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DC_VOLTAGE_DIVISOR_ID(data_ptr) \
    {                                                                                             \
        ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DC_VOLTAGE_DIVISOR_ID,                                 \
        ZB_ZCL_ATTR_TYPE_U16,                                                                     \
        ZB_ZCL_ATTR_ACCESS_READ_ONLY | ZB_ZCL_ATTR_ACCESS_REPORTING,                              \
        (ZB_ZCL_NON_MANUFACTURER_SPECIFIC),                                                       \
        (void *)data_ptr}

#define ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DC_CURRENT_DIVISOR_ID(data_ptr) \
    {                                                                                             \
        ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DC_CURRENT_DIVISOR_ID,                                 \
        ZB_ZCL_ATTR_TYPE_U16,                                                                     \
        ZB_ZCL_ATTR_ACCESS_READ_ONLY | ZB_ZCL_ATTR_ACCESS_REPORTING,                              \
        (ZB_ZCL_NON_MANUFACTURER_SPECIFIC),                                                       \
        (void *)data_ptr}

#define ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DC_POWER_DIVISOR_ID(data_ptr) \
    {                                                                                           \
        ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DC_POWER_DIVISOR_ID,                                 \
        ZB_ZCL_ATTR_TYPE_U16,                                                                   \
        ZB_ZCL_ATTR_ACCESS_READ_ONLY | ZB_ZCL_ATTR_ACCESS_REPORTING,                            \
        (ZB_ZCL_NON_MANUFACTURER_SPECIFIC),                                                     \
        (void *)data_ptr}

#define ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_ACCURRENT_DIVISOR_ID(data_ptr) \
    {                                                                                            \
        ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_ACCURRENT_DIVISOR_ID,                                 \
        ZB_ZCL_ATTR_TYPE_U16,                                                                    \
        ZB_ZCL_ATTR_ACCESS_READ_ONLY | ZB_ZCL_ATTR_ACCESS_REPORTING,                             \
        (ZB_ZCL_NON_MANUFACTURER_SPECIFIC),                                                      \
        (void *)data_ptr}

#define ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_ACCURRENT_MULTIPLIER_ID(data_ptr) \
    {                                                                                               \
        ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_ACCURRENT_MULTIPLIER_ID,                                 \
        ZB_ZCL_ATTR_TYPE_U16,                                                                       \
        ZB_ZCL_ATTR_ACCESS_READ_ONLY | ZB_ZCL_ATTR_ACCESS_REPORTING,                                \
        (ZB_ZCL_NON_MANUFACTURER_SPECIFIC),                                                         \
        (void *)data_ptr}

#define ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_ACVOLTAGE_DIVISOR_ID(data_ptr) \
    {                                                                                            \
        ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_ACVOLTAGE_DIVISOR_ID,                                 \
        ZB_ZCL_ATTR_TYPE_U16,                                                                    \
        ZB_ZCL_ATTR_ACCESS_READ_ONLY | ZB_ZCL_ATTR_ACCESS_REPORTING,                             \
        (ZB_ZCL_NON_MANUFACTURER_SPECIFIC),                                                      \
        (void *)data_ptr}

#define ZB_SET_ATTR_DESCR_WITH_ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_ACVOLTAGE_MULTIPLIER_ID(data_ptr) \
    {                                                                                               \
        ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_ACVOLTAGE_MULTIPLIER_ID,                                 \
        ZB_ZCL_ATTR_TYPE_U16,                                                                       \
        ZB_ZCL_ATTR_ACCESS_READ_ONLY | ZB_ZCL_ATTR_ACCESS_REPORTING,                                \
        (ZB_ZCL_NON_MANUFACTURER_SPECIFIC),                                                         \
        (void *)data_ptr}

// Declare this here, as the one provided in ZBOSS header is quite restricted.
// We at least want to have dc voltage, current & power. In ZBOSS there is only dc power.
#define ZB_ZCL_DECLARE_ELECTRICAL_MEASUREMENT_ATTRIB_LIST(attr_list, measurement_type, dc_voltage, dc_current, dc_power, dc_voltage_divisor, dc_current_divisor, dc_power_divisor, one, uone) \
    ZB_ZCL_START_DECLARE_ATTRIB_LIST_CLUSTER_REVISION(attr_list, ZB_ZCL_ELECTRICAL_MEASUREMENT)                                                                                               \
    /* Data attributes */                                                                                                                                                                     \
    ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_MEASUREMENT_TYPE_ID, (measurement_type))                                                                                          \
    ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DC_VOLTAGE_ID, (dc_voltage))                                                                                                      \
    ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DC_CURRENT_ID, (dc_current))                                                                                                      \
    ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DCPOWER_ID, (dc_power))                                                                                                           \
    /* Format attributes */                                                                                                                                                                   \
    ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DC_VOLTAGE_DIVISOR_ID, (dc_voltage_divisor))                                                                                      \
    ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DC_CURRENT_DIVISOR_ID, (dc_current_divisor))                                                                                      \
    ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_DC_POWER_DIVISOR_ID, (dc_power_divisor))                                                                                          \
    /* AC format attributes. Required by Z2M, as it just likes to get everything that it knows about.  */                                                                                     \
    /* https://github.com/Koenkk/zigbee-herdsman-converters/issues/11367 */                                                                                                                   \
    ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_ACTIVE_POWER_ID, (one))                                                                                                           \
    ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_ACVOLTAGE_DIVISOR_ID, (uone))                                                                                                     \
    ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_ACVOLTAGE_MULTIPLIER_ID, (uone))                                                                                                  \
    ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_ACCURRENT_DIVISOR_ID, (uone))                                                                                                     \
    ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_ACCURRENT_MULTIPLIER_ID, (uone))                                                                                                  \
    ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_ACPOWER_DIVISOR_ID, (uone))                                                                                                       \
    ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_ACPOWER_MULTIPLIER_ID, (uone))                                                                                                    \
    ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_RMSVOLTAGE_ID, (uone))                                                                                                            \
    ZB_ZCL_SET_ATTR_DESC(ZB_ZCL_ATTR_ELECTRICAL_MEASUREMENT_RMSCURRENT_ID, (uone))                                                                                                            \
    ZB_ZCL_FINISH_DECLARE_ATTRIB_LIST

typedef struct
{
    // Basic information attribute set
    zb_uint32_t measurement_type = ZB_ZCL_ELECTRICAL_MEASUREMENT_DC_MEASUREMENT;
    // DC Measuring attribute set
    zb_int16_t dc_voltage = ZB_ZCL_ELECTRICAL_MEASUREMENT_DC_VOLTAGE_DEFAULT_VALUE;
    zb_int16_t dc_current = ZB_ZCL_ELECTRICAL_MEASUREMENT_DC_CURRENT_DEFAULT_VALUE;
    zb_int16_t dc_power = ZB_ZCL_ELECTRICAL_MEASUREMENT_DCPOWER_DEFAULT_VALUE;

    // Formats
    zb_uint16_t dc_voltage_divisor = ZCL_DC_VOLTAGE_VALUE_DIVISOR;
    zb_uint16_t dc_current_divisor = ZCL_DC_CURRENT_VALUE_DIVISOR;
    zb_uint16_t dc_power_divisor = ZCL_DC_POWER_VALUE_DIVISOR;

    // Something for AC values, so Z2M can use them.
    zb_uint16_t uone = 1;
    zb_int16_t one = 1;
} zb_zcl_electrical_measurement_attrs_t;