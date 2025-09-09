#include "nrfx_temp.h"

#include "types_device_temperature.hpp"

namespace zbhome
{
    namespace components
    {
        DeviceTemperature::DeviceTemperature()
        {
            nrfx_temp_config_t config = NRFX_TEMP_DEFAULT_CONFIG;
            int err = nrfx_temp_init(&config, NULL);
            if (err != NRFX_SUCCESS)
            {
                LOG_ERR("error initing device temp sensor");
            }
        }

        void DeviceTemperature::onLoop()
        {
            int res = nrfx_temp_measure();
            if (res != NRFX_SUCCESS)
            {
                LOG_ERR("measure internal temperature failed");
                return;
            }

            int32_t device_temp = nrfx_temp_result_get();
            int32_t real_device_temp = nrfx_temp_calculate(device_temp);

            static const uint16_t clusterId = 0x0002;
            zb_zcl_status_t status = zb_zcl_set_attr_val(getEndpoint(),
                                                         clusterId,
                                                         ZB_ZCL_CLUSTER_SERVER_ROLE,
                                                         0x000, // CurrentTemperature attr. Hardcoding for now. Need better approach.
                                                         (zb_uint8_t *)&real_device_temp,
                                                         ZB_FALSE);
            if (status)
            {
                LOG_ERR("Failed to set ZCL attribute for device temp: %d", status);
            }
        }
    }
}