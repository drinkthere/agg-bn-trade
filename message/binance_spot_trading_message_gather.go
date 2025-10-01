package message

import (
	"agg-bn-trade/container"
	"agg-bn-trade/context"
	"agg-bn-trade/utils"
	"agg-bn-trade/utils/logger"
	binanceSpot "github.com/dictxwang/go-binance"
)

func StartGatherBinanceSpotOrder(orderChan chan *binanceSpot.WsUserDataEvent, globalContext *context.GlobalContext) {
	// 为每个 instID 记录上一条消息的时间戳，用于过滤
	type timeRecord struct {
		eventTime       int64
		transactionTime int64
	}
	lastTimeMap := make(map[string]*timeRecord)

	go func() {
		defer func() {
			logger.Warn("[TradingGather] Binance Spot Order Gather Exited.")
		}()
		for {
			event := <-orderChan
			order := event.OrderUpdate
			if !utils.InArray(order.Symbol, globalContext.InstrumentComposite.BinanceSpotInstIDs) {
				continue
			}

			if order.Status == string(binanceSpot.OrderStatusTypePartiallyFilled) || order.Status == string(binanceSpot.OrderStatusTypeFilled) {
				instID := order.Symbol

				// 获取该 instID 的上次时间记录
				lastTime, exists := lastTimeMap[instID]
				if exists {
					// 过滤时间戳小于上一条消息的订单
					if event.Time < lastTime.eventTime || order.TransactionTime < lastTime.transactionTime {
						logger.Debug("[TradingGather] Skip order for %s due to older timestamp: eventTime=%d < %d or transactionTime=%d < %d",
							instID, event.Time, lastTime.eventTime, order.TransactionTime, lastTime.transactionTime)
						continue
					}
				}

				globalContext.SpotOrderChannel <- &container.ZMQOrder{
					Sbl: instID,
					Px:  order.LatestPrice,
					Sz:  order.LatestVolume,
					Ets: event.Time,
					Tts: order.TransactionTime,
				}

				// 更新该 instID 的最新时间戳
				if lastTime == nil {
					lastTimeMap[instID] = &timeRecord{
						eventTime:       event.Time,
						transactionTime: order.TransactionTime,
					}
				} else {
					lastTime.eventTime = event.Time
					lastTime.transactionTime = order.TransactionTime
				}

				logger.Info("[TradingGather] instID=%s price=%s, volume=%s", instID, order.LatestPrice, order.LatestVolume)
			}
		}
	}()
	logger.Info("[TradingGather] Start Gather Binance Spot Order")
}
