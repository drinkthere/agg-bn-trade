package main

import (
	"agg-bn-trade/config"
	"agg-bn-trade/container"
	"agg-bn-trade/utils/logger"
	"encoding/json"
	zmq "github.com/pebbe/zmq4"
	"os"
	"strings"
	"time"
)

func StartZmq() {
	logger.Info("Start Binance Ticker ZMQ")
	if globalConfig.FuturesZMQIPC != "" {
		startZmq(globalConfig.FuturesZMQIPC, config.FuturesInstrument, globalContext.FuturesOrderChannel)
	}

	if globalConfig.SpotZMQIPC != "" {
		startZmq(globalConfig.SpotZMQIPC, config.SpotInstrument, globalContext.SpotOrderChannel)
	}
}

func startZmq(ipc string, instType config.InstrumentType, orderChan chan *container.ZMQOrder) {
	go func() {
		defer func() {
			logger.Warn("[StartZmq] %s %s Pub Service Listening Exited.", config.BinanceExchange, instType)
		}()

		logger.Info("[StartZmq] %s %s Start Pub Service.", config.BinanceExchange, instType)

		var ctx *zmq.Context
		var pub *zmq.Socket
		var err error
		isPubStopped := true
		for {
			if isPubStopped {
				ctx, err = zmq.NewContext()
				if err != nil {
					logger.Error("[StartZmq] %s %s New Context Error: %s", config.BinanceExchange, instType, err.Error())
					time.Sleep(time.Second * 1)
					continue
				}

				pub, err = ctx.NewSocket(zmq.PUB)
				if err != nil {
					logger.Error("[StartZmq] %s %s New Socket Error: %s", config.BinanceExchange, instType, err.Error())
					ctx.Term()
					time.Sleep(time.Second * 1)
					continue
				}

				err = pub.Bind(ipc)
				if err != nil {
					logger.Error("[StartZmq] %s %s Bind to Local ZMQ %s Error: %s", config.BinanceExchange, instType, ipc, err.Error())
					pub.Close() // 关闭套接字
					ctx.Term()  // 释放上下文
					time.Sleep(time.Second * 1)
					continue
				}

				// 如果是 IPC 文件，设置权限让其他用户可以访问
				if strings.HasPrefix(ipc, "ipc://") {
					ipcPath := strings.TrimPrefix(ipc, "ipc://")
					// 设置文件权限为 666 (rw-rw-rw-)，让所有用户都可以读写
					err := os.Chmod(ipcPath, 0666)
					if err != nil {
						logger.Warn("[StartZmq] Failed to set IPC file permissions for %s: %s", ipcPath, err.Error())
					} else {
						logger.Info("[StartZmq] Set IPC file permissions to 0666 for %s", ipcPath)
					}
				}
				logger.Info("[StartZmq] %s %s Successfully bound to %s", config.BinanceExchange, instType, ipc)
				isPubStopped = false
			}

			select {
			case zmqOrder := <-orderChan:
				jsonBytes, err := json.Marshal(zmqOrder)
				if err != nil {
					logger.Warn("[StartZmq]  %s %s  Error marshaling Ticker: %v", config.BinanceExchange, instType, err)
					continue
				}

				_, err = pub.Send(string(jsonBytes), 0)
				if err != nil {
					logger.Warn("[StartZmq] %s %s Error sending order data: %v", config.BinanceExchange, instType, err)
					isPubStopped = true
					pub.Close()
					ctx.Term()
					time.Sleep(time.Second * 5)
					continue
				}
			}
		}
	}()
}
