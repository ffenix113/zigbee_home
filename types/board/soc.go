package board

type SoC string

const (
	NRF52840 SoC = "nrf52840"
	NRF5340  SoC = "nrf5340"
	NRF54L10 SoC = "nrf54l10"
	NRF54L15 SoC = "nrf54l15"
)

var knownSoCs = map[SoC]struct{}{
	NRF52840: {},
	NRF5340:  {},
	NRF54L10: {},
	NRF54L15: {},
}

func IsKnownSoC(soc SoC) bool {
	_, ok := knownSoCs[soc]
	return ok
}
