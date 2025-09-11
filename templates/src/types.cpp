#pragma once

#include <cstdint>

#include <zephyr/device.h>
#include <zephyr/logging/log.h>

#ifdef __cplusplus
extern "C"
{
#endif

#include <zboss_api.h>
#include <zb_zcl_common.h>

#ifdef __cplusplus
}
#endif

#include "types.hpp"

LOG_MODULE_REGISTER(types, LOG_LEVEL_DBG);

namespace zbhome
{
    namespace types
    {
        void Component::setDevice(const struct device *dev)
        {
            m_device = dev;
        };

        const struct device *Component::getDevice()
        {
            return m_device;
        }

        void Component::setEndpoint(uint8_t endpoint)
        {
            m_endpoint = endpoint;
        }

        const uint8_t Component::getEndpoint()
        {
            return m_endpoint;
        }

        // Setup should do everything that is needded for component to become operational.
        // For example set default values or create some connection.
        //
        // A component may also have a specific constructor,
        // i.e. to add a device reference, or set min/max values, etc.
        bool Component::setup()
        {
            const auto *device = getDevice();
            if (device != nullptr && !device_is_ready(device))
            {
                LOG_ERR("device %s not ready", device->name);
                return false;
            }
            return true;
        };
    }
}