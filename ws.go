package main

import (
	"agg-bn-trade/config"
	"agg-bn-trade/message"
	binanceSpot "github.com/dictxwang/go-binance"
	"github.com/dictxwang/go-binance/futures"
	"time"
)

func startUserDataWebSocket() {
	futuresConfig := make([]config.BinanceConfig, 0)
	for _, cfg := range globalConfig.BinanceConfig {
		if cfg.InstrumentType == config.FuturesInstrument {
			futuresConfig = append(futuresConfig, cfg)
		}
	}

	spotConfig := make([]config.BinanceConfig, 0)
	for _, cfg := range globalConfig.BinanceConfig {
		if cfg.InstrumentType == config.SpotInstrument {
			spotConfig = append(spotConfig, cfg)
		}
	}

	if len(futuresConfig) > 0 {
		futuresOrderChan := make(chan *futures.WsUserDataEvent)
		futuresOrderLiteChan := make(chan *futures.WsUserDataEvent)
		for _, cfg := range futuresConfig {
			message.StartBinanceFuturesUserDataWs(&globalContext, cfg, futuresOrderChan, futuresOrderLiteChan)
			message.StartGatherBinanceFuturesOrder(futuresOrderChan, futuresOrderLiteChan, &globalContext)
			time.Sleep(1 * time.Second)
		}
	}

	if len(spotConfig) > 0 {
		spotOrderChan := make(chan *binanceSpot.WsUserDataEvent)
		for _, cfg := range spotConfig {
			message.StartBinanceSpotUserDataWs(&globalContext, cfg, spotOrderChan)
			message.StartGatherBinanceSpotOrder(spotOrderChan, &globalContext)
			time.Sleep(1 * time.Second)
		}
	}

}
