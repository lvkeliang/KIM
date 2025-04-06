package services

import (
	"KIM/communication/tcp"
	"KIM/config"
	"KIM/container"
	"KIM/logger"
	"KIM/naming"
	"KIM/naming/consul"
	"KIM/protocol"
	"KIM/services"
	"KIM/services/services/conf"
	"KIM/services/services/login"
	"KIM/storage/redis"
	"context"
)

type ServerStartOptions struct {
	Config      string
	ServiceName string
	RunType     string
}

func RunServerStart(ctx context.Context, opts *ServerStartOptions, version string) error {
	initConf, err := conf.Init(opts.Config)
	if err != nil {
		return err
	}
	_ = logger.Init(logger.Settings{
		Level: "trace",
	})

	// 指令路由
	r := services.NewRouter()
	// login
	loginHandler := login.NewLoginHandler()
	r.Handle(protocol.CommandLoginSignIn, loginHandler.DoSysLogin)
	r.Handle(protocol.CommandLoginSignOut, loginHandler.DoSysLogout)

	rdb, err := conf.InitRedis(initConf.RedisAddrs, "")
	if err != nil {
		return err
	}

	// 会话管理
	cache := redis.NewRedisStorage(rdb)
	servHandler := NewServHandler(r, cache)

	service := &naming.DefaultServiceRegistration{
		Id:       initConf.ServiceID,
		Name:     opts.ServiceName,
		Address:  initConf.PublicAddress,
		Port:     initConf.PublicPort,
		Protocol: string(protocol.ProtocolTCP),
		Tags:     initConf.Tags,
	}

	if opts.RunType == "local" {
		service.Address = "127.0.0.1"
	}

	srv := tcp.NewServer(initConf.Listen, service)

	srv.SetReadWait(config.DefaultReadWait)
	srv.SetAcceptor(servHandler)
	srv.SetMessageListener(servHandler)
	srv.SetStateListener(servHandler)

	if err := container.Init(srv); err != nil {
		return err
	}

	//ns, err := etcd.NewEtcdNaming([]string{initConf.ConsulURL})
	ns, err := consul.NewNaming(initConf.ConsulURL)
	if err != nil {
		return err
	}
	container.SetServiceNaming(ns)

	return container.Start()
}
