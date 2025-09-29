package client

import (
	"agg-bn-trade/config"
	"github.com/dictxwang/go-binance/futures"
)

type BinanceClient struct {
	FuturesClient   *futures.Client
	FuturesWsClient *futures.ClientWs
	limitProcess    int
}

func (cli *BinanceClient) Init(cfg *config.BinanceConfig) bool {
	if cfg.Intranet {
		// 如果配置的是内网，这里需要设置一下
		futures.UseIntranet = true
	}
	if cfg.LocalIP == "" {
		cli.FuturesClient = futures.NewClient(cfg.APIKey, cfg.APISecret)
	} else {
		cli.FuturesClient = futures.NewClientWithIP(cfg.APIKey, cfg.APISecret, cfg.LocalIP)
	}
	return true
}
