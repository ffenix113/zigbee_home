#pragma once

#include <cstdint>

#include <zephyr/device.h>

#ifdef __cplusplus
extern "C"
{
#endif

#include <zboss_api.h>
#include <zb_zcl_common.h>

#ifdef __cplusplus
}
#endif

namespace zbhome
{
    namespace types
    {
        class Component
        {
        public:
            // If the component requires some device - this would be the function to add it.
            // It will be added if component(==Sensor in Go code)
            // definition mentiones that it uses device.
            void setDevice(const struct device *dev);
            const struct device *getDevice();

            // setEndpoint is separate function so that implementations
            // would not need to define `setup` function and call
            // parent's `setup`.
            // Instead this function will always be called by codegen.
            void setEndpoint(uint8_t endpoint);

            const uint8_t getEndpoint();

            // Setup should do everything that is needded for component to become operational.
            // For example set default values or create some connection.
            //
            // A component may also have a specific constructor,
            // i.e. to add a device reference, or set min/max values, etc.
            virtual bool setup();

        private:
            uint8_t m_endpoint = 0;
            const struct device *m_device = nullptr;
        };

        // Other components inherit from Component because we needed access
        // to getEndpoint() from ZCLCommandHandler, and as ZCLCommandHandler
        // would only have zcl handler function - it was not possible.
        // It should not result in different resulting class size,
        // but would allow access to Component functions.
        //
        // For now it works, but maybe better separation would be prefered.

        class Sensor : public Component
        {
        public:
            // onLoop will be called on each iteration.
            //
            // No reason to return error here. If to break loop - other sensors may not run.
            // Maybe only do some metrics in the future?..
            // Let per-sensor logger(if any) log.
            virtual void onLoop() = 0;
        };

        class ZCLCommandHandler : public Component
        {
        public:
            // I would much rather prefer to have args like (cluster, attr, value),
            // but the value inside does not define a type..
            // Open to suggestions to improve/resolve this.
            virtual void zclSetAttrValue(zb_zcl_set_attr_value_param_t *setValueParam) = 0;
        };
    }
}