package config

type (
	Exchange       string
	InstrumentType string
)

const (
	BinanceExchange = Exchange("Binance")
	OkxExchange     = Exchange("Okx")

	UnknownInstrument = InstrumentType("UNKNOWN")
	SpotInstrument    = InstrumentType("SPOT")
	FuturesInstrument = InstrumentType("FUTURES")
)
