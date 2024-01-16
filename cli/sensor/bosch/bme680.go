package bosch

type BME680 struct {
	BME280
}

func NewBME680() *BME680 {
	return &BME680{
		BME280: BME280{
			variant: "bme680",
		},
	}
}

func (BME680) String() string {
	return "Bosch BME680"
}
