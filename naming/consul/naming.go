package consul

import (
	"KIM/inter"
	"KIM/logger"
	"KIM/naming"
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"sync"
	"time"
)

const (
	KeyProtocol  = "protocol"
	KeyHealthURL = "health_url"
)

// Watch 监控一个服务
type Watch struct {
	Service   string
	Callback  func([]inter.ServiceRegistration)
	WaitIndex uint64
	Quit      chan struct{}
}

type Naming struct {
	sync.RWMutex
	cli     *api.Client
	watches map[string]*Watch
}

func (n *Naming) Find(serviceName string, tags ...string) ([]inter.ServiceRegistration, error) {
	services, _, err := n.load(serviceName, 0, tags...)
	if err != nil {
		return nil, err
	}
	return services, nil
}

// waitIndex: 是否阻塞查询，0为不阻塞
func (n *Naming) load(name string, waitIndex uint64, tags ...string) ([]inter.ServiceRegistration, *api.QueryMeta, error) {
	opts := &api.QueryOptions{
		UseCache:  true,
		MaxAge:    time.Minute,
		WaitIndex: waitIndex,
	}
	catalogServices, meta, err := n.cli.Catalog().ServiceMultipleTags(name, tags, opts)
	if err != nil {
		return nil, meta, err
	}
	services := make([]inter.ServiceRegistration, 0, len(catalogServices))
	for _, s := range catalogServices {
		if s.Checks.AggregatedStatus() != api.HealthPassing {
			logger.Debugf("load service: id:%s name:%s %s:%d Status:%s", s.ServiceID, s.ServiceName, s.ServiceAddress, s.ServicePort, s.Checks.AggregatedStatus())
			continue
		}
		services = append(services, &naming.DefaultServiceRegistration{
			Id:       s.ServiceID,
			Name:     s.ServiceName,
			Address:  s.ServiceAddress,
			Port:     s.ServicePort,
			Protocol: s.ServiceMeta[KeyProtocol],
			Tags:     s.ServiceTags,
			Meta:     s.ServiceMeta,
		})
	}
	logger.Debugf("load service: %v, meta: %v", services, meta)
	return services, meta, nil
}

func (n *Naming) Subscribe(serviceName string, callback func(services []inter.ServiceRegistration)) error {
	n.Lock()
	defer n.Unlock()
	if _, ok := n.watches[serviceName]; ok {
		return errors.New("serviceName has already been registered")
	}

	w := &Watch{
		Service:  serviceName,
		Callback: callback,
		Quit:     make(chan struct{}, 1),
	}

	n.watches[serviceName] = w
	go n.watch(w)
	return nil
}

func (n *Naming) watch(wh *Watch) {
	stopped := false
	var doWatch = func(service string, callback func([]inter.ServiceRegistration)) {
		services, meta, err := n.load(service, wh.WaitIndex)
		if err != nil {
			logger.Warn(err)
			return
		}
		select {
		// 由于n.load会阻塞，所以要等待n.load返回后才会结束
		case <-wh.Quit:
			stopped = true
			logger.Infof("watch %s stopped", wh.Service)
			return
		default:
		}
		wh.WaitIndex = meta.LastIndex
		if callback != nil {
			callback(services)
		}
	}

	// 初始化 WaitIndex
	doWatch(wh.Service, nil)

	// 由于n.load会阻塞，所以要等待n.load返回后才会结束
	for !stopped {
		//这之后的doWatch中的n.load由于wh.WaitIndex已经有值了，所以会阻塞等待数据变化才返回
		doWatch(wh.Service, wh.Callback)
	}
}

func (n *Naming) Unsubscribe(serviceName string) error {
	n.Lock()
	defer n.Unlock()
	wh, ok := n.watches[serviceName]

	delete(n.watches, serviceName)
	if ok {
		close(wh.Quit)
	}
	return nil
}

func (n *Naming) Register(service inter.ServiceRegistration) error {
	reg := &api.AgentServiceRegistration{
		ID:      service.ServiceID(),
		Name:    service.ServiceName(),
		Tags:    service.GetTags(),
		Port:    service.PublicPort(),
		Address: service.PublicAddress(),
		Meta:    service.GetMeta(),
		// consul的免费版本不支持namespace
	}

	if reg.Meta == nil {
		reg.Meta = make(map[string]string)
	}

	// 把service.GetProtocol()添加到Meta中，让服务消费方知道服务提供的接入协议
	reg.Meta[KeyProtocol] = service.GetProtocol()

	//consul健康检查
	healthURL := service.GetMeta()[KeyHealthURL]
	if healthURL != "" {
		check := new(api.AgentServiceCheck)
		check.CheckID = fmt.Sprintf("%s_normal", service.ServiceID())
		check.HTTP = healthURL
		check.Timeout = "1s" // http timeout
		check.Interval = "10s"
		check.DeregisterCriticalServiceAfter = "20s" // 服务故障20s后，Agent会把它下线
		reg.Check = check
	}

	err := n.cli.Agent().ServiceRegister(reg)
	return err
}

func (n *Naming) Deregister(serviceID string) error {
	return n.cli.Agent().ServiceDeregister(serviceID)
}

func NewNaming(consulUrl string) (naming.Naming, error) {
	conf := api.DefaultConfig()
	conf.Address = consulUrl
	cli, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}
	newNaming := &Naming{
		cli:     cli,
		watches: make(map[string]*Watch, 1),
	}

	return newNaming, nil
}
