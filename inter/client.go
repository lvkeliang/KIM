package inter

import (
	"KIM/protocol"
	"net"
	"time"
)

type Client interface {
	ID() string
	Name() string
	// Connect 向一个服务器地址发起连接
	Connect(string) error
	// SetDialer 设置拨号器, Dialer在Connect中会被调用，用于完成连接的建立和握手
	SetDialer(Dialer)
	// Send 发送消息到服务端
	Send([]byte) error
	// Read 读取一帧数据, 底层复用Conn, 所以直接返回Frame
	Read() (protocol.Frame, error)
	Close()
}

// Dialer 在Connect中会被调用，用于完成连接的建立和握手
type Dialer interface {
	DialAndHandshake(DialerContext) (net.Conn, error)
}

type DialerContext struct {
	Id      string
	Name    string
	Address string
	Timeout time.Duration
}
