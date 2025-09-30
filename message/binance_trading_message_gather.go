package message

import (
	"agg-bn-trade/container"
	"agg-bn-trade/context"
	"agg-bn-trade/utils"
	"agg-bn-trade/utils/logger"
	"github.com/dictxwang/go-binance/futures"
)

func StartGatherBinanceFuturesOrder(orderChan chan *futures.WsUserDataEvent, orderLiteChan chan *futures.WsUserDataEvent, globalContext *context.GlobalContext) {
	// 为每个 instID 记录上一条消息的时间戳，用于过滤
	type timeRecord struct {
		eventTime       int64
		transactionTime int64
	}
	lastTimeMap := make(map[string]*timeRecord)

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

				// 获取该 instID 的上次时间记录
				lastTime, exists := lastTimeMap[instID]
				if exists {
					// 过滤时间戳小于上一条消息的订单
					if event.Time < lastTime.eventTime || event.TransactionTime < lastTime.transactionTime {
						logger.Debug("[TradingGather] Skip order for %s due to older timestamp: eventTime=%d < %d or transactionTime=%d < %d",
							instID, event.Time, lastTime.eventTime, event.TransactionTime, lastTime.transactionTime)
						continue
					}
				}

				globalContext.FuturesOrderChannel <- &container.ZMQOrder{
					Sbl: instID,
					Px:  order.OriginalPrice,
					Sz:  order.LastFilledQty,
					Ets: event.Time,
					Tts: event.TransactionTime,
				}

				// 更新该 instID 的最新时间戳
				if lastTime == nil {
					lastTimeMap[instID] = &timeRecord{
						eventTime:       event.Time,
						transactionTime: event.TransactionTime,
					}
				} else {
					lastTime.eventTime = event.Time
					lastTime.transactionTime = event.TransactionTime
				}

				logger.Info("[TradingGather] instID=%s price=%s, volume=%s", instID, order.OriginalPrice, order.LastFilledQty)
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
			// 获取该 instID 的上次时间记录
			lastTime, exists := lastTimeMap[instID]
			if exists {
				// 过滤时间戳小于上一条消息的订单
				if order.Time < lastTime.eventTime || order.TransactionTime < lastTime.transactionTime {
					logger.Debug("[TradingGatherLite] Skip order for %s due to older timestamp: eventTime=%d < %d or transactionTime=%d < %d",
						instID, order.Time, lastTime.eventTime, order.TransactionTime, lastTime.transactionTime)
					continue
				}
			}

			globalContext.FuturesOrderChannel <- &container.ZMQOrder{
				Sbl: instID,
				Px:  order.LastFilledPrice,
				Sz:  order.LastFilledQuantity,
				Ets: order.Time,
				Tts: order.TransactionTime,
			}

			// 更新该 instID 的最新时间戳
			if lastTime == nil {
				lastTimeMap[instID] = &timeRecord{
					eventTime:       order.Time,
					transactionTime: order.TransactionTime,
				}
			} else {
				lastTime.eventTime = order.Time
				lastTime.transactionTime = order.TransactionTime
			}

			logger.Info("[TradingGatherLite] instID=%s price=%f, volume=%f", instID, order.LastFilledPrice, order.LastFilledQuantity)
		}
	}()
	logger.Info("[TradingGatherLite] Start Gather Binance Futures Lite Order")
}
