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

#ifdef __cplusplus
extern "C"
{
#endif

// #include <zboss_api_buf.h>
#include <zboss_api.h>
#include <zb_zcl_el_measurement.h>

#ifdef __cplusplus
}
#endif

#include "cluster_electrical_measurement.hpp"

#ifdef __cplusplus
extern "C"
{
#endif

    static zb_ret_t check_value_el_measurement_server(zb_uint16_t attr_id, zb_uint8_t endpoint, zb_uint8_t *value)
    {
        ZVUNUSED(attr_id);
        ZVUNUSED(value);
        ZVUNUSED(endpoint);

        /* All values for mandatory attributes are allowed, extra check for
         * optional attributes is needed */

        return RET_OK;
    }

    void zb_zcl_el_measurement_init_server()
    {
        zb_zcl_add_cluster_handlers(ZB_ZCL_CLUSTER_ID_ELECTRICAL_MEASUREMENT,
                                    ZB_ZCL_CLUSTER_SERVER_ROLE,
                                    (zb_zcl_cluster_check_value_t)check_value_el_measurement_server,
                                    (zb_zcl_cluster_write_attr_hook_t)NULL,
                                    (zb_zcl_cluster_handler_t)NULL);
    }

    void zb_zcl_el_measurement_init_client()
    {
        zb_zcl_add_cluster_handlers(ZB_ZCL_CLUSTER_ID_ELECTRICAL_MEASUREMENT,
                                    ZB_ZCL_CLUSTER_CLIENT_ROLE,
                                    (zb_zcl_cluster_check_value_t)NULL,
                                    (zb_zcl_cluster_write_attr_hook_t)NULL,
                                    (zb_zcl_cluster_handler_t)NULL);
    }

#ifdef __cplusplus
}
#endif