package main

import (
	"KIM/services/gateway"
	"context"
	"flag"
)

var (
	configPath = flag.String("c", "conf.yaml", "配置文件路径")
	protocol   = flag.String("p", "ws", "服务协议 (ws/tcp)")
	route      = flag.String("r", "/", "WebSocket路由路径")
	version    = "v1.0.0" // 可通过编译时注入
)

func main() {
	flag.Parse()

	opts := &gateway.ServerStartOptions{
		Config:   *configPath,
		Protocol: *protocol,
		Route:    *route,
	}

	if err := gateway.RunServerStart(
		context.Background(),
		opts,
		version,
	); err != nil {
		panic(err) // 实际项目应使用日志记录
	}
}
