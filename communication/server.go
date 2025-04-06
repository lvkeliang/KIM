package communication

import (
	"KIM/protocol"
	"context"
	"time"
)

// Service 定义了基础服务的抽象接口
type Service interface {
	ServiceID() string
	ServiceName() string
	GetMeta() map[string]string
}

// ServiceRegistration 服务注册接口
type ServiceRegistration interface {
	Service
	GetTags() []string
	PublicPort() int
	DialURL() string
	PublicAddress() string
	GetProtocol() string
}

type Server interface {
	ServiceRegistration
	SetAcceptor(Acceptor)
	SetMessageListener(MessageListener)
	SetStateListener(StateListener)
	// SetReadWait 设置连接超时, 控制心跳逻辑
	SetReadWait(duration time.Duration)
	// SetChannelMap 设置连接管理器，自动管理连接的生命周期
	SetChannelMap(ChannelMap)
	Start() error
	Push(channelId string, data []byte) error
	Shutdown(ctx context.Context) error
}

// Acceptor 处理握手
type Acceptor interface {
	// Accept 返回channelID和error
	Accept(protocol.Conn, time.Duration) (string, error)
}

// StateListener 状态监听器, 上报连接事件给业务层, 让业务层实现一些处理逻辑
type StateListener interface {
	Disconnect(string) error
}

// MessageListener 消息监听器
type MessageListener interface {
	// Receive 第一个参数Agent表示发送方
	// 第二个参数[]byte 必须是一个或多个完整数据包, 否则上层业务会拆包失败
	Receive(Agent, []byte)
}

type Agent interface {
	// ID 返回连接的channelID
	ID() string
	// Push 用于上层业务返回消息, 要求线程安全
	Push([]byte) error
}
