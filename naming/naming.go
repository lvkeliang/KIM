package naming

import "KIM/inter"

// Naming 服务发现
type Naming interface {
	Find(serviceName string, tags ...string) ([]inter.ServiceRegistration, error)            //服务发现，支持通过tag查询
	Subscribe(serviceName string, callback func(services []inter.ServiceRegistration)) error //订阅服务变更通知
	Unsubscribe(serviceName string) error                                                    //取消订阅服务变更通知
	Register(service inter.ServiceRegistration) error                                        //注册服务
	Deregister(serviceID string) error                                                       //注销服务
}
