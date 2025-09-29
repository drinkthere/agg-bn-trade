package config

import (
	"encoding/json"
	"go.uber.org/zap/zapcore"
	"os"
)

type BinanceConfig struct {
	Account        string
	APIKey         string
	APISecret      string
	LocalIP        string
	Intranet       bool
	InstrumentType InstrumentType
}

type Config struct {
	// 日志配置
	LogLevel zapcore.Level
	LogPath  string

	BinanceConfig         []BinanceConfig
	BinanceFuturesInstIDs []string // 需要收集的InstID
	BinanceSpotInstIDs    []string // 需要收集的InstID

	MinAccuracy float64 // 价格最小精度

	FuturesZMQIPC string
	SpotZMQIPC    string
}

func LoadConfig(filename string) *Config {
	config := new(Config)
	reader, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	// 加载配置
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}
