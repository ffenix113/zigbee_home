package board

var knownSoCs = map[string]struct{}{
	"nrf52840": {},
	"nrf5340":  {},
	"nrf54l05": {},
	"nrf54l10": {},
	"nrf54l15": {},
}

func IsKnownSoC(soc string) bool {
	_, ok := knownSoCs[soc]
	return ok
}
