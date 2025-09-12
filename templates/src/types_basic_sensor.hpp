#pragma once

#include <cstdint>

#include "types.hpp"

namespace zbhome
{
    namespace components
    {
        // BasicSensor will update and set all available
        // sensor channels. Currently it is not possible to
        // limit channels that should be checked.
        class BasicSensor : public zbhome::types::Sensor
        {
        public:
            void onLoop();
        };
    };
};