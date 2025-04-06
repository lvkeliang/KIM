package naming

import (
	"KIM/communication"
)

// Naming 服务发现
type Naming interface {
	Find(serviceName string, tags ...string) ([]communication.ServiceRegistration, error)            //服务发现，支持通过tag查询
	Subscribe(serviceName string, callback func(services []communication.ServiceRegistration)) error //订阅服务变更通知
	Unsubscribe(serviceName string) error                                                            //取消订阅服务变更通知
	Register(service communication.ServiceRegistration) error                                        //注册服务
	Deregister(serviceID string) error                                                               //注销服务
}
