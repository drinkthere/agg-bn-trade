package context

import (
	"agg-bn-trade/config"
	"agg-bn-trade/container"
)

type GlobalContext struct {
	InstrumentComposite *container.InstrumentComposite
	FuturesOrderChannel chan *container.ZMQOrder
	SpotOrderChannel    chan *container.ZMQOrder
}

func (context *GlobalContext) Init(globalConfig *config.Config) {

	// 初始化交易对数据
	context.initInstrumentComposite(globalConfig)

	context.FuturesOrderChannel = make(chan *container.ZMQOrder)
	context.SpotOrderChannel = make(chan *container.ZMQOrder)
}

func (context *GlobalContext) initInstrumentComposite(globalConfig *config.Config) {
	instrumentComposite := container.NewInstrumentComposite(globalConfig)
	context.InstrumentComposite = instrumentComposite
}
