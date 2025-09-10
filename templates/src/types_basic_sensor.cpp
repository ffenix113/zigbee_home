#include <cstdint>

#include <zephyr/drivers/sensor.h>
#include <zephyr/logging/log.h>

#include "types_basic_sensor.hpp"
#include "zbhome_sensor.hpp"

LOG_MODULE_REGISTER(basic_sensor, LOG_LEVEL_INF);

namespace zbhome
{
    namespace components
    {
        void BasicSensor::onLoop()
        {
            const auto sensorDevice = getDevice();

            sensor_sample_fetch(sensorDevice);

            zbhome::sensors::updateFetchedSamples(sensorDevice, getEndpoint());
        }
    };
}