package cluster

type Temperature struct{}

func (t Temperature) ID() ID {
	return ID_TEMP_MEASUREMENT
}
