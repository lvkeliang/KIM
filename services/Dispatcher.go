package services

import (
	"KIM/container"
	"KIM/protocol"
	"strings"
)

// Dispatcher 消息分发器，用于推送消息到网关
type Dispatcher interface {
	Push(gateway string, channels []string, packet *protocol.LogicPkt) error // 向网关gateway中的channels每个连接推送一条LogicPkt消息
}

type DispatcherImpl struct {
}

func (d *DispatcherImpl) Push(gateway string, channels []string, packet *protocol.LogicPkt) error {
	packet.AddStringMeta(protocol.MetaDestChannels, strings.Join(channels, ","))
	return container.Push(gateway, packet)
}
