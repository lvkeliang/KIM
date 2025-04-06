package gateway

import (
	"KIM/communication"
	"KIM/communication/websocket"
	"KIM/container"
	"KIM/logger"
	"KIM/naming"
	"KIM/naming/consul"
	"KIM/protocol"
	"KIM/services/gateway/conf"
	"KIM/services/gateway/serv"
	"context"
	"time"
)

type ServerStartOptions struct {
	Config   string
	Protocol string
	Route    string
}

func RunServerStart(ctx context.Context, opts *ServerStartOptions, version string) error {
	config, err := conf.Init(opts.Config)
	if err != nil {
		return err
	}
	_ = logger.Init(logger.Settings{
		Level: "trace",
	})

	handler := &serv.Handler{ServiceId: config.ServiceID}

	var srv communication.Server
	seivice := &naming.DefaultServiceRegistration{
		Id:       config.ServiceID,
		Name:     config.ServiceName,
		Address:  config.PublicAddress,
		Port:     config.PublicPort,
		Protocol: opts.Protocol,
		Tags:     config.Tags,
	}
	if opts.Protocol == "ws" {
		srv = websocket.NewServer(config.Listen, seivice)
	}

	srv.SetReadWait(time.Minute * 2)
	srv.SetAcceptor(handler)
	srv.SetMessageListener(handler)
	srv.SetStateListener(handler) // 添加这行

	_ = container.Init(srv, protocol.SNChat, protocol.SNLogin)

	//ns, err := etcd.NewEtcdNaming([]string{config.ConsulURL})
	ns, err := consul.NewNaming(config.ConsulURL)
	if err != nil {
		return err
	}
	container.SetServiceNaming(ns)

	container.SetDialer(serv.NewDialer(config.ServiceID))

	return container.Start()
}
