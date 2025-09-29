package message

import (
	"agg-bn-trade/container"
	"agg-bn-trade/context"
	"agg-bn-trade/utils"
	"agg-bn-trade/utils/logger"
	"github.com/dictxwang/go-binance/futures"
)

func StartGatherBinanceFuturesOrder(orderChan chan *futures.WsUserDataEvent, orderLiteChan chan *futures.WsUserDataEvent, globalContext *context.GlobalContext) {

	go func() {
		defer func() {
			logger.Warn("[TradingGather] Binance Futures Order Gather Exited.")
		}()
		for {
			event := <-orderChan
			order := event.OrderTradeUpdate
			if !utils.InArray(order.Symbol, globalContext.InstrumentComposite.BinanceFuturesInstIDs) {
				continue
			}

			if order.Status == futures.OrderStatusTypePartiallyFilled || order.Status == futures.OrderStatusTypeFilled {
				instID := order.Symbol
				globalContext.FuturesOrderChannel <- &container.ZMQOrder{
					Sbl: instID,
					Px:  order.OriginalPrice,
					Sz:  order.LastFilledQty,
					Ets: event.Time,
					Tts: event.TransactionTime,
				}
				logger.Debug("[TradingGather] instID=%s price=%s, volume=%s", instID, order.OriginalPrice, order.LastFilledQty)
			}
		}
	}()
	logger.Info("[TradingGather] Start Gather Binance Futures Order")

	go func() {
		defer func() {
			logger.Warn("[TradingGatherLite] Binance Futures Order Gather Exited.")
		}()
		for {
			order := <-orderLiteChan
			if !utils.InArray(order.Symbol, globalContext.InstrumentComposite.BinanceFuturesInstIDs) {
				continue
			}

			instID := order.Symbol
			globalContext.FuturesOrderChannel <- &container.ZMQOrder{
				Sbl: instID,
				Px:  order.LastFilledPrice,
				Sz:  order.LastFilledQuantity,
				Ets: order.Time,
				Tts: order.TransactionTime,
			}
			logger.Debug("[TradingGatherLite] instID=%s price=%s, volume=%s", instID, order.LastFilledPrice, order.LastFilledQuantity)
		}
	}()
	logger.Info("[TradingGatherLite] Start Gather Binance Futures Lite Order")
}
