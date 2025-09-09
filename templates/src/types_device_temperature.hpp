#pragma once

#include "types.hpp"

namespace zbhome
{
    namespace components
    {
        class DeviceTemperature : public zbhome::types::Sensor
        {
        public:
            DeviceTemperature();
            void onLoop();
        };
    }
}