package container

import (
	"KIM/config"
	"KIM/inter"
	"KIM/logger"
	"KIM/naming"
	"KIM/protocol"
	"KIM/protocol/protoImpl"
	"KIM/tcp"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	stateUninitialized = iota
	stateInitialized
	stateStarted
	stateClosed
)

const (
	StateYoung = "young"
	StateAdult = "adult"
)

const (
	KeyServiceState = "service_state"
)

type Container struct {
	sync.RWMutex
	naming.Naming
	Srv        inter.Server
	state      uint32
	srvclients map[string]ClientMap
	selector   Selector
	dialer     inter.Dialer        //创建Client时用于处理拨号和握手
	deps       map[string]struct{} // 依赖的服务
}

var log = logger.WithField("module", "container")

var c = &Container{
	state:    0,
	selector: &HashSelector{},
	deps:     make(map[string]struct{}),
}

func Default() *Container {
	return c
}

// Init deps参数传入依赖的服务，如wire.SNChat, wire.SNLogin
func Init(srv inter.Server, deps ...string) error {
	if !atomic.CompareAndSwapUint32(&c.state, stateUninitialized, stateInitialized) {
		return errors.New("container has already initialized")
	}
	c.Srv = srv
	for _, dep := range deps {
		if _, ok := c.deps[dep]; ok {
			continue
		}
		c.deps[dep] = struct{}{}
	}
	log.WithField("func", "Init").Infof("sr %s:%s - deps %v", srv.ServiceID(), srv.ServiceName(), c.deps)
	c.srvclients = make(map[string]ClientMap, len(deps))
	return nil
}

// SetDialer 创建Client时用于处理拨号和握手
func SetDialer(dialer inter.Dialer) {
	c.dialer = dialer
}

// SetSelector 注册一个自定义的服务路由器
func SetSelector(selector Selector) {
	c.selector = selector
}

func Start() error {
	if c.Naming == nil {
		return fmt.Errorf("naming is nil")
	}

	if !atomic.CompareAndSwapUint32(&c.state, stateInitialized, stateStarted) {
		return errors.New("container has started")
	}

	// 启动server
	go func(srv inter.Server) {
		err := srv.Start()
		if err != nil {
			log.Errorln(err)
		}
	}(c.Srv)

	// 与依赖的服务连接
	for service := range c.deps {
		go func(service string) {
			err := connectToService(service)
			if err != nil {
				log.Errorln(err)
			}
		}(service)
	}

	// 服务注册
	if c.Srv.PublicAddress() != "" && c.Srv.PublicPort() != 0 {
		err := c.Naming.Register(c.Srv)
		if err != nil {
			log.Errorln(err)
		}
	}

	// 阻塞并等待退出
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	log.Infoln("shutdown", <-c)
	// 退出
	return shutdown()
}

func connectToService(serviceName string) error {
	clients := NewClients(10)
	c.srvclients[serviceName] = clients

	//1.订阅服务的新增
	delay := time.Second * 10
	err := c.Naming.Subscribe(serviceName, func(services []inter.ServiceRegistration) {
		for _, service := range services {
			if _, ok := clients.Get(service.ServiceID()); ok {
				continue
			}
			log.WithField("func", "connectToService").Infof("Watch a new service: %v", service)
			service.GetMeta()[KeyServiceState] = StateYoung
			go func(service inter.ServiceRegistration) {
				// 等待新增的服务与各个服务完成连接
				time.Sleep(delay)
				service.GetMeta()[KeyServiceState] = StateAdult
			}(service)

			_, err := buildClient(clients, service)
			if err != nil {
				logger.Warn(err)
			}
		}
	})

	if err != nil {
		return err
	}

	//2.查询已经存在的服务
	services, err := c.Naming.Find(serviceName)
	if err != nil {
		return err
	}
	log.Info("find service ", services)
	for _, service := range services {
		// 标记为StateAdult
		service.GetMeta()[KeyServiceState] = StateAdult
		_, err := buildClient(clients, service)
		if err != nil {
			logger.Warn(err)
		}
	}
	return nil
}

// 发现提供服务的server后，立即创建一个client与之建立连接
func buildClient(clients ClientMap, service inter.ServiceRegistration) (inter.Client, error) {
	c.Lock()
	defer c.Unlock()
	var (
		id   = service.ServiceID()
		name = service.ServiceName()
		meta = service.GetMeta()
	)
	// 1. 检测连接是否已经存在
	if _, ok := clients.Get(id); ok {
		return nil, nil
	}
	// 2. 服务之间只允许使用tcp协议
	if service.GetProtocol() != string(protocol.ProtocolTCP) {
		return nil, fmt.Errorf("unexpected service Protocol: %s", service.GetProtocol())
	}

	// 3. 构建客户端并建立连接
	cli := tcp.NewClientWithProps(id, name, meta, tcp.ClientOptions{
		Heartbeat: config.DefaultHeartbeat,
		ReadWait:  config.DefaultReadWait,
		WriteWait: config.DefaultWriteWait,
	})
	if c.dialer == nil {
		return nil, fmt.Errorf("dialer is nil")
	}
	cli.SetDialer(c.dialer)
	err := cli.Connect(service.DialURL())
	if err != nil {
		return nil, err
	}
	// 4. 读取消息
	go func(cli inter.Client) {
		err := readLoop(cli)
		if err != nil {
			log.Debug(err)
		}
		clients.Remove(id)
		cli.Close()
	}(cli)
	// 5. 添加到客户端集合中
	clients.Add(cli)
	return cli, nil
}

// Receive default listener
func readLoop(cli inter.Client) error {
	log := logger.WithFields(logger.Fields{
		"module": "container",
		"func":   "readLoop",
	})
	log.Infof("readLoop started of %s %s", cli.ServiceID(), cli.ServiceName())
	for {
		frame, err := cli.Read()
		if err != nil {
			log.Trace(err)
			return err
		}
		if frame.GetOpCode() != protocol.OpBinary {
			continue
		}
		buf := bytes.NewBuffer(frame.GetPayload())

		packet, err := protocol.MustUnMarshalLogicPkt(buf)
		if err != nil {
			log.Info(err)
			continue
		}
		err = pushMessage(packet)
		if err != nil {
			log.Info(err)
		}
	}
}

// 消息通过网关服务器推送到channel中
func pushMessage(packet *protocol.LogicPkt) error {
	server, _ := packet.GetSpecificMeta(protocol.MetaDestServer)
	if server != c.Srv.ServiceID() {
		return fmt.Errorf("meta dest_server is incorrect, %s != %s", server, c.Srv.ServiceID())
	}
	channels, ok := packet.GetSpecificMeta(protocol.MetaDestChannels)
	if !ok {
		return fmt.Errorf("meta dest_channels is nil")
	}

	channelIds := strings.Split(channels.(string), ",")
	packet.DelMeta(protocol.MetaDestServer)
	packet.DelMeta(protocol.MetaDestChannels)
	payload := protocol.MarshalPacket(packet)
	log.Debugf("Push to %v %v", channelIds, packet)

	for _, channel := range channelIds {
		//性能监测网关下发字节数
		//TODO: messageOutFlowBytes.WithLabelValues(packet.Command).Add(float64(len(payload)))
		err := c.Srv.Push(channel, payload)
		if err != nil {
			log.Debug(err)
		}
	}
	return nil
}

func shutdown() error {
	if !atomic.CompareAndSwapUint32(&c.state, stateStarted, stateClosed) {
		return errors.New("container has closed")
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()

	//关闭服务器
	err := c.Srv.Shutdown(ctx)
	if err != nil {
		log.Errorln(err)
	}

	//从注册中心注销服务
	err = c.Naming.Deregister(c.Srv.ServiceID())
	if err != nil {
		log.Warn(err)
	}

	// 退订服务变更
	for dep := range c.deps {
		_ = c.Naming.Unsubscribe(dep)
	}

	log.Infoln("shutdown")
	return nil
}

// Push 消息下行
// 通过创建并设置Client的Dialer和Server的Accepter，握手时同步Client创建的ID给Server，Server将此ID作为ChannelID存到ChannelMap
func Push(channelID string, p *protocol.LogicPkt) error {
	p.AddStringMeta(protocol.MetaDestServer, channelID)
	return c.Srv.Push(channelID, protocol.MarshalPacket(p))
}

// Forward 消息上行
func Forward(serviceName string, packet *protocol.LogicPkt) error {
	if packet == nil {
		return errors.New("packet is nil")
	}
	if packet.Command == "" {
		return errors.New("command is empty in packet")
	}
	if packet.ChannelId == "" {
		return errors.New("channelId is empty in packet")
	}
	return ForwardWithSelector(serviceName, packet, c.selector)
}

func ForwardWithSelector(serviceName string, packet *protocol.LogicPkt, selector Selector) error {
	cli, err := lookup(serviceName, &packet.Header, selector)
	if err != nil {
		return err
	}

	//add a tag in packet
	packet.AddStringMeta(protocol.MetaDestServer, c.Srv.ServiceID())
	log.Debugf("forward message to %v with %s", cli.ServiceID(), &packet.Header)
	return cli.Send(protocol.MarshalPacket(packet))
}

// 负载均衡和路由，根据服务名查找一个可靠的服务
func lookup(serviceName string, header *protoImpl.Header, selector Selector) (inter.Client, error) {
	clients, ok := c.srvclients[serviceName]
	if !ok {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}

	// 只获取状态为StateAdult的服务
	srvs := clients.Services(KeyServiceState, StateAdult)
	if len(srvs) == 0 {
		return nil, fmt.Errorf("no services found for %s", serviceName)
	}
	id := selector.Lookup(header, srvs)
	if cli, ok := clients.Get(id); ok {
		return cli, nil
	}
	return nil, fmt.Errorf("no client found")
}
