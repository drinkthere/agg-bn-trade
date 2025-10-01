package message

import (
	"agg-bn-trade/client"
	"agg-bn-trade/config"
	mmcontext "agg-bn-trade/context"
	"agg-bn-trade/utils/logger"
	"context"
	binanceSpot "github.com/dictxwang/go-binance"
	"time"
)

func StartBinanceSpotUserDataWs(globalContext *mmcontext.GlobalContext,
	binanceConfig config.BinanceConfig, orderChan chan *binanceSpot.WsUserDataEvent) {
	orderWs := binanceSpotOrderWebSocket{
		cfg:       binanceConfig,
		orderChan: orderChan,
		isStopped: true,
	}

	orderWs.keepAlive()
	orderWs.startSpotUserDataWs(globalContext)
	logger.Info("[OrderWebSocket] Start Listen Binance Spot Order")
}

type binanceSpotOrderWebSocket struct {
	cfg       config.BinanceConfig
	listenKey string
	orderChan chan *binanceSpot.WsUserDataEvent
	isStopped bool
	stopChan  chan struct{}
}

func (ws *binanceSpotOrderWebSocket) keepAlive() {
	var binanceClient client.BinanceClient
	binanceClient.Init(&ws.cfg)

	// 获取 listenKey，监听transaction 消息时，需要这个 key
	listenKey, err := binanceClient.SpotClient.NewStartMarginUserStreamService().Do(context.Background())
	if err != nil {
		logger.Error("[BinanceUserDataWs] %s Get binance spot listen key failed, exit the program", ws.cfg.Account)
		return
	}
	ws.listenKey = listenKey

	// listenKey 每60分钟过期一次，所以需要加个定时器，提前续期
	go func() {
		for {
			time.Sleep(30 * time.Minute)
			err = binanceClient.SpotClient.NewKeepaliveMarginUserStreamService().ListenKey(ws.listenKey).Do(context.Background())
			if err != nil {
				logger.Error("[BinanceUserDataWs] %s Refresh binance spot listen key failed, exit the program", ws.cfg.Account)
				return
			}
		}
	}()
}

func (ws *binanceSpotOrderWebSocket) handleUserDataEvent(event *binanceSpot.WsUserDataEvent) {
	if event.Event == binanceSpot.UserDataEventTypeExecutionReport {
		ws.orderChan <- event
	}
}

func (ws *binanceSpotOrderWebSocket) handleError(err error) {
	// 出错断开连接，再重连
	logger.Error("[BinanceUserDataWs] %s Binance SPOT User Data Handle Error And Will Reconnect Ws in 1 Seconds: %s", ws.cfg.Account, err.Error())
	ws.stopChan <- struct{}{}
	ws.isStopped = true
}

func (ws *binanceSpotOrderWebSocket) startSpotUserDataWs(globalContext *mmcontext.GlobalContext) {
	go func() {
		defer func() {
			logger.Warn("[BinanceUserDataWs] %s Binance Spot User Data Listening Exited.", ws.cfg.Account)
		}()
		for {
			if !ws.isStopped {
				time.Sleep(time.Second * 1)
				continue
			}

			var stopChan chan struct{}
			var err error

			if ws.cfg.LocalIP == "" {
				_, stopChan, err = binanceSpot.WsUserDataServe(ws.listenKey, ws.handleUserDataEvent, ws.handleError)
			} else {
				_, stopChan, err = binanceSpot.WsUserDataServeWithIp(ws.cfg.LocalIP, ws.listenKey, ws.handleUserDataEvent, ws.handleError)
			}
			if err != nil {
				logger.Error("[BinanceUserDataWs] %s Subscribe Binance Spot User Data Error: %s", ws.cfg.Account, err.Error())
				time.Sleep(time.Second * 1)
				continue
			}
			logger.Info("[BinanceUserDataWs] %s Subscribe Binance Spot User Data Success", ws.cfg.Account)
			// 重置channel和时间
			ws.stopChan = stopChan
			ws.isStopped = false
		}
	}()
}
