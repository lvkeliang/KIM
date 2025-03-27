package main

import (
	"KIM/inter"
	"KIM/logger"
	"KIM/protocol"
	"KIM/tcp"
	"KIM/websocket"
	"context"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
	"strings"
	"time"
)

func main() {

	client := &ClientDemo{}
	client.Start("minimal_client", "tcp", "http://127.0.0.1:8080")
	// client.Start("minimal_client", "ws", "ws://127.0.0.1:8080")

}

type ClientDemo struct{}

func (c *ClientDemo) Start(userID, protocolwt, addr string) {
	// 初始化客户端
	var cli inter.Client
	if protocolwt == "ws" {
		cli = websocket.NewClient(userID, "client", websocket.ClientOptions{})
		cli.SetDialer(&WebsocketDialer{})
	} else if protocolwt == "tcp" {
		cli = tcp.NewClient("test1", "client", tcp.ClientOptions{})
		cli.SetDialer(&TCPDialer{})
	}
	// 建立连接
	err := cli.Connect(addr)
	if err != nil {
		logger.Error(err)
	}
	count := 10
	go func() {
		// 发送消息然后退出
		for i := 0; i < count; i++ {
			err := cli.Send([]byte("hello"))
			if err != nil {
				logger.Error(err)
				return
			}
			time.Sleep(time.Second)
		}
	}()

	// 接收消息
	recv := 0
	for {
		frame, err := cli.Read()
		if err != nil {
			logger.Info(err)
			break
		}
		if frame.GetOpCode() != protocol.OpBinary {
			continue
		}
		recv++
		logger.Warnf("%s receive message [%s]", cli.ID(), frame.GetPayload())

		if recv == count { // 接收完消息
			break
		}
	}

	// 退出
	cli.Close()
}

type ClientHandler struct{}

func (h *ClientHandler) Receive(ag inter.Agent, payload []byte) {
	logger.Warnf("%s receive message [%s]", ag.ID, string(payload))
}

func (h *ClientHandler) Disconnect(id *string) error {
	logger.Warnf("disconnect %s", id)
	return nil
}

type WebsocketDialer struct {
	userID string
}

func (d *WebsocketDialer) DialAndHandshake(ctx inter.DialerContext) (net.Conn, error) {
	// 调用ws.Dial拨号
	conn, _, _, err := ws.Dial(context.TODO(), ctx.Address)
	if err != nil {
		return nil, err
	}
	// 发送认证消息, 示例为userid
	err = wsutil.WriteClientBinary(conn, []byte(ctx.Id))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// TCPDialer TCPDialer
type TCPDialer struct {
}

// DialAndHandshake DialAndHandshake
func (d *TCPDialer) DialAndHandshake(ctx inter.DialerContext) (net.Conn, error) {
	logger.Info("start tcp dial: ", ctx.Address)
	// 1 调用net.Dial拨号
	address := ctx.Address
	if strings.Contains(address, "://") {
		// 去除协议前缀
		address = strings.Split(address, "://")[1]
	}
	conn, err := net.DialTimeout("tcp", address, ctx.Timeout)
	if err != nil {
		return nil, err
	}
	// 2. 发送用户认证信息，示例就是userid
	err = tcp.WriteFrame(conn, protocol.OpBinary, []byte(ctx.Id))
	if err != nil {
		return nil, err
	}
	// 3. return conn
	return conn, nil
}
