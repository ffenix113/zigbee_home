package bosch

func NewBME680() *BME280 {
	return &BME280{
		Variant: "bme680",
	}
}
