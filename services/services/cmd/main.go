package main

import (
	"KIM/protocol"
	"KIM/services/services"
	"context"
	"flag"
)

var (
	configPath  = flag.String("c", "conf.yaml", "配置文件路径")
	serviceName = flag.String("n", protocol.SNLogin, "服务注册名称")
	version     = "v1.0.0" // 可通过编译时注入
)

func main() {
	flag.Parse()

	opts := &services.ServerStartOptions{
		Config:      *configPath,
		ServiceName: *serviceName,
		RunType:     "local",
	}

	if err := services.RunServerStart(
		context.Background(),
		opts,
		version,
	); err != nil {
		panic(err) // 实际项目应使用日志记录
	}
}
