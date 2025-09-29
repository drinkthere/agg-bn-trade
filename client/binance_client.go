package client

import (
	"agg-bn-trade/config"
	"context"
	"github.com/dictxwang/go-binance/futures"
	"github.com/juju/ratelimit"
)

type BinanceClient struct {
	FuturesClient   *futures.Client
	FuturesWsClient *futures.ClientWs

	bucket10s    *ratelimit.Bucket
	bucket60s    *ratelimit.Bucket
	limitProcess int
}

func (cli *BinanceClient) Init(cfg *config.BinanceConfig) bool {
	if cfg.Intranet {
		// 如果配置的是内网，这里需要设置一下
		futures.UseIntranet = true
	}

	ctx := context.Background()
	if cfg.LocalIP == "" {
		cli.FuturesClient = futures.NewClient(cfg.APIKey, cfg.APISecret)
		cli.FuturesWsClient = futures.NewTradingWsClient(ctx, cfg.APIKey, cfg.APISecret, "")
	} else {
		cli.FuturesClient = futures.NewClientWithIP(cfg.APIKey, cfg.APISecret, cfg.LocalIP)
		cli.FuturesWsClient = futures.NewTradingWsClient(ctx, cfg.APIKey, cfg.APISecret, cfg.LocalIP)
	}
	return true
}
