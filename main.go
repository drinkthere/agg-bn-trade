package main

import (
	"agg-bn-trade/config"
	"agg-bn-trade/context"
	"agg-bn-trade/utils"
	"agg-bn-trade/utils/logger"
	"fmt"
	_ "net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
)

var globalConfig config.Config
var globalContext context.GlobalContext

func ExitProcess() {
	os.Exit(1)
}

func main() {
	runtime.GOMAXPROCS(1)
	// 参数判断
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s config_file\n", os.Args[0])
		os.Exit(1)
	}

	// 监听退出消息，并调用ExitProcess进行处理
	utils.RegisterExitSignal(ExitProcess)

	// 加载配置文件
	globalConfig = *config.LoadConfig(os.Args[1])

	// 设置日志级别, 并初始化日志
	logger.InitLogger(globalConfig.LogPath, globalConfig.LogLevel)

	// 解析config，加载杠杆和合约交易对，初始化context，账户初始化设置，拉取仓位、余额等
	globalContext.Init(&globalConfig)

	// 开始监听order消息
	startUserDataWebSocket()

	// 启动zmq发送成交的order消息
	StartZmq()

	// 阻塞主进程
	for {
		time.Sleep(24 * time.Hour)
	}
}
