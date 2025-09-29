package container

import (
	"agg-bn-trade/config"
)

type InstrumentComposite struct {
	BinanceFuturesInstIDs []string
	BinanceSpotInstIDs    []string
}

func NewInstrumentComposite(globalConfig *config.Config) *InstrumentComposite {
	composite := &InstrumentComposite{
		BinanceFuturesInstIDs: globalConfig.BinanceFuturesInstIDs,
		BinanceSpotInstIDs:    globalConfig.BinanceSpotInstIDs,
	}

	return composite
}
