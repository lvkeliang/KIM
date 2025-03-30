package container

import (
	"KIM/inter"
	"KIM/protocol/protoImpl"
)

// Selector 用于在消息上行时，从一批服务列表中选择一个合适的服务
type Selector interface {
	Lookup(*protoImpl.Header, []inter.Service) string
}
