package inter

import (
	"KIM/protocol"
	"time"
)

type Channel interface {
	protocol.Conn
	Agent
	Close() error                       // 重写net.Conn中的Close方法
	ReadLoop(lst MessageListener) error // 阻塞的方法, 把消息的读取和心跳处理封装在一起
	SetWriteWait(time.Duration)
	SetReadWait(time.Duration)
}

// ChannelMap 连接管理器, 管理Channel
type ChannelMap interface {
	Add(channel Channel)
	Remove(id string)
	Get(id string) (Channel, bool)
	All() []Channel
}
