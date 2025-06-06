package tcp

import (
	"KIM/communication"
	"KIM/config"
	"KIM/logger"
	"KIM/protocol"
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ClientOptions ClientOptions
type ClientOptions struct {
	Heartbeat time.Duration //登录超时
	ReadWait  time.Duration //读超时
	WriteWait time.Duration //写超时
}

// Client is a tcp implement of kim.Client
type Client struct {
	sync.Mutex
	communication.Dialer
	once    sync.Once
	id      string
	name    string
	conn    protocol.Conn
	state   int32
	options ClientOptions
	Meta    map[string]string
}

func (c *Client) ServiceID() string {
	return c.id
}

func (c *Client) ServiceName() string {
	return c.name
}

func (c *Client) GetMeta() map[string]string {
	return c.Meta
}

// Connect to server
func (c *Client) Connect(addr string) error {
	// 如果 addr 是 URL（如 "ws://..."），则解析；否则直接作为 "host:port"
	var hostPort string
	if strings.Contains(addr, "://") {
		u, err := url.Parse(addr)
		if err != nil {
			return fmt.Errorf("invalid URL: %v", err)
		}
		hostPort = u.Host // 提取 "host:port"
	} else {
		hostPort = addr // 直接使用 "host:port"
	}

	// 验证 host:port 格式
	if _, _, err := net.SplitHostPort(hostPort); err != nil {
		return fmt.Errorf("invalid address format: %v", err)
	}

	// 这里是一个CAS原子操作，对比并设置值，是并发安全的。
	if !atomic.CompareAndSwapInt32(&c.state, 0, 1) {
		return fmt.Errorf("client has connected")
	}

	rawconn, err := c.Dialer.DialAndHandshake(communication.DialerContext{
		Id:      c.id,
		Name:    c.name,
		Address: addr,
		Timeout: config.DefaultLoginWait,
	})
	if err != nil {
		atomic.CompareAndSwapInt32(&c.state, 1, 0)
		return err
	}
	if rawconn == nil {
		return fmt.Errorf("conn is nil")
	}
	c.conn = NewConn(rawconn)

	if c.options.Heartbeat > 0 {
		go func() {
			err := c.heartbealoop()
			if err != nil {
				logger.WithField("module", "tcp.client").Warn("heartbealoop stopped - ", err)
			}
		}()
	}
	return nil
}

func (c *Client) Send(payload []byte) error {
	if atomic.LoadInt32(&c.state) == 0 {
		return fmt.Errorf("connection is nil")
	}
	c.Lock()
	defer c.Unlock()
	err := c.conn.SetWriteDeadline(time.Now().Add(c.options.WriteWait))
	if err != nil {
		return err
	}
	return c.conn.WriteFrame(protocol.OpBinary, payload)
}

func (c *Client) Read() (protocol.Frame, error) {
	if c.conn == nil {
		return nil, errors.New("connection is nil")
	}
	if c.options.Heartbeat > 0 {
		_ = c.conn.SetReadDeadline(time.Now().Add(c.options.ReadWait))
	}
	frame, err := c.conn.ReadFrame()
	if err != nil {
		return nil, err
	}
	if frame.GetOpCode() == protocol.OpClose {
		return nil, errors.New("remote side close the channel")
	}
	return frame, nil
}

func (c *Client) heartbealoop() error {
	tick := time.NewTicker(c.options.Heartbeat)
	for range tick.C {
		// 发送一个ping的心跳包给服务端
		if err := c.ping(); err != nil {
			return err
		}
	}
	return nil
}

// SetDialer 设置握手逻辑
func (c *Client) SetDialer(dialer communication.Dialer) {
	c.Dialer = dialer
}

func (c *Client) Close() {
	c.once.Do(func() {
		if c.conn == nil {
			return
		}
		// graceful close connection
		_ = wsutil.WriteClientMessage(c.conn, ws.OpClose, nil)

		c.conn.Close()
		atomic.CompareAndSwapInt32(&c.state, 1, 0)
	})
}

func (c *Client) ping() error {
	logger.WithField("module", "tcp.client").Tracef("%s send ping to server", c.id)

	err := c.conn.SetWriteDeadline(time.Now().Add(c.options.WriteWait))
	if err != nil {
		return err
	}
	return c.conn.WriteFrame(protocol.OpPing, nil)
}

// NewClient NewClient
func NewClient(id, name string, opts ClientOptions) communication.Client {
	if opts.WriteWait == 0 {
		opts.WriteWait = config.DefaultWriteWait
	}
	if opts.ReadWait == 0 {
		opts.ReadWait = config.DefaultReadWait
	}
	cli := &Client{
		id:      id,
		name:    name,
		options: opts,
	}
	return cli
}

func NewClientWithProps(id, name string, meta map[string]string, opts ClientOptions) communication.Client {
	if opts.WriteWait == 0 {
		opts.WriteWait = config.DefaultWriteWait
	}
	if opts.ReadWait == 0 {
		opts.ReadWait = config.DefaultReadWait
	}

	cli := &Client{
		id:      id,
		name:    name,
		options: opts,
		Meta:    meta,
	}
	return cli
}
