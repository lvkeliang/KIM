package tcp

import (
	"KIM/config"
	"KIM/inter"
	"KIM/logger"
	"KIM/protocol"
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/ksuid"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type defaultAcceptor struct {
}

// NewDefaultAcceptor 创建默认Acceptor
func NewDefaultAcceptor() *defaultAcceptor {
	return &defaultAcceptor{}
}

// Accept 实现Acceptor接口
func (a *defaultAcceptor) Accept(conn protocol.Conn, timeout time.Duration) (string, error) {
	return ksuid.New().String(), nil
}

type ServerOptions struct {
	loginWait time.Duration //登录超时
	readWait  time.Duration //读超时
	writeWait time.Duration //写超时
}

// Server is a tcp implement of kim.Server
type Server struct {
	listen string
	inter.ServiceRegistration
	inter.ChannelMap
	inter.Acceptor
	inter.MessageListener
	inter.StateListener
	once    sync.Once
	options ServerOptions
	quit    int32
}

// Start server
func (s *Server) Start() error {
	log := logger.WithFields(logger.Fields{
		"module": "tcp.server",
		"listen": s.listen,
		"id":     s.ServiceID(),
	})

	if s.Acceptor == nil {
		s.Acceptor = new(defaultAcceptor)
	}

	if s.StateListener == nil {
		return fmt.Errorf("StateListener is nil")
	}

	if s.ChannelMap == nil {
		s.ChannelMap = NewChannels(100)
	}
	// step 1
	lst, err := net.Listen("tcp", s.listen)
	if err != nil {
		return err
	}
	log.Info("started")
	for {
		// step 2
		rawconn, err := lst.Accept()
		if err != nil {
			rawconn.Close()
			log.Warn(err)
			continue
		}
		go func(rawconn net.Conn) {
			conn := NewConn(rawconn)
			// step 3
			id, err := s.Accept(conn, s.options.loginWait)
			if err != nil {
				_ = conn.WriteFrame(protocol.OpClose, []byte(err.Error()))
				conn.Close()
				return
			}

			if _, ok := s.Get(id); ok {
				log.Warnf("channel %s existed", id)
				_ = conn.WriteFrame(protocol.OpClose, []byte("channelId is repeated"))
				conn.Close()
				return
			}
			//step 4
			channel := NewChannel(id, conn)
			channel.SetReadWait(s.options.readWait)
			channel.SetWriteWait(s.options.writeWait)
			s.Add(channel)

			log.Info("accept ", channel.ID())
			//step 5
			err = channel.ReadLoop(s.MessageListener)
			if err != nil {
				//TODO: 处理TCP正常退出的unexpected EOF错误
				log.Info(err)
			}
			// step 6
			s.Remove(channel.ID())
			_ = s.Disconnect(channel.ID())
			channel.Close()
		}(rawconn)
	}
}

func (s *Server) SetAcceptor(acceptor inter.Acceptor) {
	s.Acceptor = acceptor
}

func (s *Server) SetMessageListener(messageListener inter.MessageListener) {
	s.MessageListener = messageListener
}

func (s *Server) SetStateListener(stateListener inter.StateListener) {
	s.StateListener = stateListener
}

func (s *Server) SetReadWait(duration time.Duration) {
	s.options.readWait = duration
}

func (s *Server) SetChannelMap(channelMap inter.ChannelMap) {
	s.ChannelMap = channelMap
}

// string channelID
// []byte data
func (s *Server) Push(id string, data []byte) error {
	ch, ok := s.ChannelMap.Get(id)
	if !ok {
		return errors.New("channel no found")
	}
	return ch.Push(data)
}

func (s *Server) Shutdown(ctx context.Context) error {
	log := logger.WithFields(logger.Fields{
		// "module": s.Name(),
		"id": s.ServiceID(),
	})
	s.once.Do(func() {
		defer func() {
			log.Infoln("shutdown")
		}()
		if atomic.CompareAndSwapInt32(&s.quit, 0, 1) {
			return
		}

		// close channels
		chanels := s.ChannelMap.All()
		for _, ch := range chanels {
			ch.Close()

			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
	})
	return nil
}

func NewServer(listen string, service inter.ServiceRegistration) inter.Server {
	return &Server{
		listen:              listen,
		ServiceRegistration: service,
		quit:                0,
		options: ServerOptions{
			loginWait: config.DefaultLoginWait,
			readWait:  config.DefaultReadWait,
			writeWait: config.DefaultWriteWait,
		},
	}
}
