
namespace zbhome {

namespace types {

class Device {
public:
    // Setup should do everything that is needded for device to become operational.
    // For example set default values or create some connection.
    virtual void setup() {};
}

class Sensor {
public:
    // on_loop will be called on each iteration.
    //
    // No reason to return error here. Let logger(if any) log.
    virtual void on_loop() = 0;
};


}