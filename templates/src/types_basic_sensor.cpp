#include <cstdint>

#include <zephyr/drivers/sensor.h>

#include "types.hpp"
#include "types_basic_sensor.hpp"
#include "zbhome_sensor.hpp"

namespace zbhome
{

    namespace types
    {
        // sensorType will define which channels
        BasicSensor::BasicSensor(SensorType sensorType) : m_sensorType(sensorType)
        {
        }

        void BasicSensor::onLoop()
        {
            // This is something that should not happen.
            if (this->m_sensorType == 0)
            {
                return;
            }

            const auto sensorDevice = getDevice();
            sensor_sample_fetch(sensorDevice);

            // Technically we can get all supported channels here
            // and use them without the need for specifying sensor type.
            // This would require to compare the return value to -ENOTSUP.

            if (m_sensorType & Temperature)
            {
                zbhome_sensor_fetch_and_update_temperature(sensorDevice, getEndpoint());
            }

            if (m_sensorType & Humidity)
            {
                zbhome_sensor_fetch_and_update_humidity(sensorDevice, getEndpoint());
            }

            if (m_sensorType & Pressure)
            {
                zbhome_sensor_fetch_and_update_pressure(sensorDevice, getEndpoint());
            }

            #ifdef ZB_ZCL_CLUSTER_ID_CARBON_DIOXIDE
            if (m_sensorType & CarbonDioxide)
            {
                zbhome_sensor_fetch_and_update_carbon_dioxide(sensorDevice, getEndpoint());
            }
            #endif
        }
    };
}