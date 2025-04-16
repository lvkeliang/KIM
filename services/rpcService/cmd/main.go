package main

import (
	"KIM/services/rpcService"
	"context"
	"flag"
)

var (
	configPath = flag.String("c", "conf.yaml", "配置文件路径")
	version    = "v1.0.0" // 可通过编译时注入
)

func main() {
	flag.Parse()

	opts := &rpcService.ServerStartOptions{
		Config:  *configPath,
		RunType: "docker",
	}

	if err := rpcService.RunServerStart(
		context.Background(),
		opts,
		version,
	); err != nil {
		panic(err) // 实际项目应使用日志记录
	}
}
