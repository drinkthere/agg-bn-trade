package container

import (
	"agg-bn-trade/config"
)

type ZMQOrder struct {
	Sbl string
	Px  string
	Sz  string
	Ets int64
	Tts int64
}

type Order struct {
	Exchange      config.Exchange
	InstID        string
	OrderPrice    float64
	OrderVolume   float64 // Okx合约是张数
	EventTs       int64
	TransactionTs int64
}
