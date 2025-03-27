package mock

import (
	"KIM/inter"
	"KIM/logger"
	"KIM/protocol"
	"KIM/tcp"
	"KIM/websocket"
	"errors"
	"time"
)

type ServerDemo struct{}

type ServerHandler struct{}

func (s *ServerDemo) Start(id, protocol, addr string) {
	var srv inter.Server
	service := inter.ServiceRegistration{}

	if protocol == "ws" {
		srv = websocket.NewServer(addr, service)
	} else if protocol == "tcp" {
		srv = tcp.NewServer(addr, service)
	}

	handler := &ServerHandler{}

	srv.SetReadWait(time.Minute)
	srv.SetAcceptor(handler)
	srv.SetMessageListener(handler)

	err := srv.Start()
	if err != nil {
		panic(err)
	}
}

func (h *ServerHandler) Accept(conn protocol.Conn, timeout time.Duration) (string, error) {
	// 读取客户端发送的鉴权数据包
	frame, err := conn.ReadFrame()
	if err != nil {
		return "", err
	}
	// 解析数据包, 数据包内容就是userid
	userID := string(frame.GetPayload())
	// 鉴权
	if userID == "" {
		return "", errors.New("user id is invalid")
	}
	return userID, nil
}

func (h *ServerHandler) Receive(ag inter.Agent, payload []byte) {
	ack := string(payload) + " from server "
	_ = ag.Push([]byte(ack))
}

func (h *ServerHandler) Disconnect(id string) error {
	logger.Warnf("disconnect %s", id)
	return nil
}
