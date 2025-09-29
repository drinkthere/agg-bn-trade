package main

import (
	"agg-bn-trade/config"
	"agg-bn-trade/message"
	"github.com/dictxwang/go-binance/futures"
	"time"
)

func startUserDataWebSocket() {
	// 监听Binance统一账户永续订单信息 和持仓信息
	futuresOrderChan := make(chan *futures.WsUserDataEvent)
	futuresOrderLiteChan := make(chan *futures.WsUserDataEvent)

	futuresConfig := make([]config.BinanceConfig, 0)
	for _, cfg := range globalConfig.BinanceConfig {
		if cfg.InstrumentType == config.FuturesInstrument {
			futuresConfig = append(futuresConfig, cfg)
		}
	}

	for _, cfg := range futuresConfig {
		message.StartBinanceFuturesUserDataWs(&globalContext, cfg, futuresOrderChan, futuresOrderLiteChan)
		message.StartGatherBinanceFuturesOrder(futuresOrderChan, futuresOrderLiteChan, &globalContext)
		time.Sleep(1 * time.Second)
	}
}
