#pragma once

#include <cstdint>

#include <zephyr/device.h>

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
            void setDevice(const struct device* dev) {
                m_device = dev;
            }

            const struct device* getDevice() {
                return m_device;
            }

            // setEndpoint is separate function so that implementations
            // would not need to define `setup` function and call
            // parent's `setup`.
            // Instead this function will always be called by codegen.
            void setEndpoint(uint8_t endpoint) {
                m_endpoint = endpoint;
            }

            uint8_t getEndpoint()
            {
                return m_endpoint;
            }

            // Setup should do everything that is needded for component to become operational.
            // For example set default values or create some connection.
            //
            // A component may also have a specific constructor,
            // i.e. to add a device reference, or set min/max values, etc.
            virtual bool setup() {
                return true;
            };

        private:
            uint8_t m_endpoint = 0;
            const struct device *m_device = nullptr;
        };

        class Sensor
        {
        public:
            // onLoop will be called on each iteration.
            //
            // No reason to return error here. If to break loop - other sensors may not run.
            // Maybe only do some metrics in the future?..
            // Let per-sensor logger(if any) log.
            virtual void onLoop() = 0;
        };
    }
}