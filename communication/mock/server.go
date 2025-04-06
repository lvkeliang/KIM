package main

import (
	"KIM/communication"
	"KIM/communication/tcp"
	"KIM/communication/websocket"
	"KIM/logger"
	"KIM/naming"
	"KIM/protocol"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// 创建服务器实例
	server := &ServerDemo{}

	// 启动服务器
	go server.Start("minimal_server", "tcp", ":8080")
	//go server.Start("minimal_server", "ws", ":8080")

	// 等待终止信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	logger.Info("server shutting down...")
}

type ServerDemo struct{}

type ServerHandler struct{}

func (s *ServerDemo) Start(id, protocol, addr string) {
	var srv communication.Server
	service := naming.DefaultServiceRegistration{
		Id:       id,
		Protocol: protocol,
	}

	if protocol == "ws" {
		srv = websocket.NewServer(addr, service)
	} else if protocol == "tcp" {
		srv = tcp.NewServer(addr, service)
	}

	handler := &ServerHandler{}

	// 必须设置所有必需的监听器
	srv.SetReadWait(time.Minute)
	srv.SetAcceptor(handler)
	srv.SetMessageListener(handler)

	srv.SetStateListener(handler) // 添加这行

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

func (h *ServerHandler) Receive(ag communication.Agent, payload []byte) {
	ack := string(payload) + " from server"
	_ = ag.Push([]byte(ack))
}

func (h *ServerHandler) Disconnect(id string) error {
	logger.Warnf("disconnect %s", id)
	return nil
}
