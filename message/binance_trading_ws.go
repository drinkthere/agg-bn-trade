package message

import (
	"agg-bn-trade/client"
	"agg-bn-trade/config"
	mmcontext "agg-bn-trade/context"
	"agg-bn-trade/utils/logger"
	"context"
	"github.com/dictxwang/go-binance/futures"
	"time"
)

func StartBinanceFuturesUserDataWs(globalContext *mmcontext.GlobalContext,
	binanceConfig config.BinanceConfig, orderChan chan *futures.WsUserDataEvent, orderLiteChan chan *futures.WsUserDataEvent) {
	orderWs := binanceFuturesOrderWebSocket{
		cfg:           binanceConfig,
		orderChan:     orderChan,
		orderLiteChan: orderLiteChan,
		isStopped:     true,
	}

	orderWs.keepAlive()
	orderWs.startFuturesUserDataWs(globalContext)
	logger.Info("[OrderWebSocket] Start Listen Binance Futures Order")
}

type binanceFuturesOrderWebSocket struct {
	cfg           config.BinanceConfig
	listenKey     string
	orderChan     chan *futures.WsUserDataEvent
	orderLiteChan chan *futures.WsUserDataEvent
	isStopped     bool
	stopChan      chan struct{}
}

func (ws *binanceFuturesOrderWebSocket) keepAlive() {
	var binanceClient client.BinanceClient
	binanceClient.Init(&ws.cfg)

	// 获取 listenKey，监听transaction 消息时，需要这个 key
	listenKey, err := binanceClient.FuturesClient.NewStartUserStreamService().Do(context.Background())
	if err != nil {
		logger.Error("[BinanceUserDataWs] %s Get binance futures listen key failed, exit the program", ws.cfg.Account)
		return
	}
	ws.listenKey = listenKey

	// listenKey 每60分钟过期一次，所以需要加个定时器，提前续期
	go func() {
		for {
			time.Sleep(30 * time.Minute)
			err = binanceClient.FuturesClient.NewKeepaliveUserStreamService().ListenKey(ws.listenKey).Do(context.Background())
			if err != nil {
				logger.Error("[BinanceUserDataWs] %s Refresh binance futures listen key failed, exit the program", ws.cfg.Account)
				return
			}
		}
	}()
}

func (ws *binanceFuturesOrderWebSocket) handleUserDataEvent(event *futures.WsUserDataEvent) {
	if event.Event == futures.UserDataEventTypeOrderTradeUpdate {
		ws.orderChan <- event
	} else if event.Event == futures.UserDataEventTradeLite {
		ws.orderLiteChan <- event
	}
}

func (ws *binanceFuturesOrderWebSocket) handleError(err error) {
	// 出错断开连接，再重连
	logger.Error("[BinanceUserDataWs] %s Binance Futures User Data Handle Error And Will Reconnect Ws in 1 Seconds: %s", ws.cfg.Account, err.Error())
	ws.stopChan <- struct{}{}
	ws.isStopped = true
}

func (ws *binanceFuturesOrderWebSocket) startFuturesUserDataWs(globalContext *mmcontext.GlobalContext) {
	go func() {
		defer func() {
			logger.Warn("[BinanceUserDataWs] %s Binance Futures User Data Listening Exited.", ws.cfg.Account)
		}()
		for {
			if !ws.isStopped {
				time.Sleep(time.Second * 1)
				continue
			}

			var stopChan chan struct{}
			var err error

			futures.UseIntranet = ws.cfg.Intranet
			if ws.cfg.LocalIP == "" {
				_, stopChan, err = futures.WsUserDataServe(ws.listenKey, ws.handleUserDataEvent, ws.handleError)
			} else {
				_, stopChan, err = futures.WsUserDataServeWithIP(ws.cfg.LocalIP, ws.listenKey, ws.handleUserDataEvent, ws.handleError)
			}
			if err != nil {
				logger.Error("[BinanceUserDataWs] %s Subscribe Binance Futures User Data Error: %s", ws.cfg.Account, err.Error())
				time.Sleep(time.Second * 1)
				continue
			}
			logger.Info("[BinanceUserDataWs] %s Subscribe Binance Futures User Data Success", ws.cfg.Account)
			// 重置channel和时间
			ws.stopChan = stopChan
			ws.isStopped = false
		}
	}()
}
