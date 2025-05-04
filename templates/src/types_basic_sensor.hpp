#pragma once

#include <cstdint>

#include "types.hpp"

namespace zbhome {
    namespace types {
        // If there is easier way to do this, like enums - that would be so much better.
        // Open to suggestions/solutions!
        typedef uint8_t SensorType;
        constexpr SensorType Temperature    = 1 << 1;
        constexpr SensorType Humidity       = 1 << 2;
        constexpr SensorType Pressure       = 1 << 3;
        constexpr SensorType CarbonDioxide  = 1 << 4;
        
        class BasicSensor : public Component {
            public:
                BasicSensor(SensorType sensorType);
                void onLoop();
            private:
                SensorType m_sensorType = 0;
        };
    };
};