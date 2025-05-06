#pragma once

#include <cstdint>

#include "types.hpp"

namespace zbhome {
    namespace types {
        // BasicSensor will update and set all available
        // sensor channels. Currently it is not possible to 
        // limit channels that should be checked.
        class BasicSensor : public Component, public Sensor {
            public:
                void onLoop();
        };
    };
};