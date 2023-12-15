package sensor

type Sensor interface {
}

// Sensor defines all information necessary about the attached sensor.
type Base struct {
	Type string
	// Connection provides information about communication protocol.
	Connection map[string]string
}
