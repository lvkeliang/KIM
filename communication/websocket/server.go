package websocket

import (
	"KIM/communication"
	"KIM/config"
	"KIM/logger"
	"KIM/protocol"
	"context"
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/segmentio/ksuid"
	"net/http"
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
	//// 1. 设置读取超时
	//_ = conn.SetReadDeadline(time.Now().Add(timeout))
	//defer conn.SetReadDeadline(time.Time{}) // 重置超时
	//
	//// 2. 读取握手消息
	//frame, err := conn.ReadFrame()
	//if err != nil {
	//	return "", errors.New("read handshake frame failed: " + err.Error())
	//}
	//
	//// 3. 验证消息类型
	//if frame.GetOpCode() != protocol.OpBinary {
	//	return "", errors.New("invalid handshake opcode")
	//}
	//
	//// 4. 解析握手数据（示例JSON格式）
	//var handshake struct {
	//	Token    string `json:"token"`
	//	DeviceID string `json:"device_id"`
	//}
	//if err := json.Unmarshal(frame.GetPayload(), &handshake); err != nil {
	//	return "", errors.New("invalid handshake format")
	//}
	//
	//// 5. 验证Token（示例逻辑）
	//if handshake.Token == "" {
	//	return "", errors.New("token required")
	//}
	//
	//// 6. 生成ChannelID（示例：deviceID + timestamp）
	//if handshake.DeviceID == "" {
	//	return "", errors.New("device_id required")
	//}
	//channelID := handshake.DeviceID + "-" + time.Now().Format("20060102150405")
	//
	//// 7. 返回成功响应
	//response := map[string]string{"status": "success", "channel_id": channelID}
	//respData, _ := json.Marshal(response)
	//if err := conn.WriteFrame(protocol.OpBinary, respData); err != nil {
	//	return "", errors.New("send handshake response failed")
	//}
	//
	//return channelID, nil
}

type ServerOptions struct {
	loginWait time.Duration //登录超时
	readWait  time.Duration //读超时
	writeWait time.Duration //写超时
}

type Server struct {
	listen string
	communication.ServiceRegistration
	communication.ChannelMap
	communication.Acceptor
	communication.MessageListener
	communication.StateListener
	once    sync.Once
	options ServerOptions
	quit    int32
}

// resp 发送HTTP响应
func resp(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(message))
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	log := logger.WithFields(logger.Fields{
		"module": "ws.server",
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
		s.ChannelMap = communication.NewChannels(100)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rawconn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			resp(w, http.StatusBadRequest, err.Error())
			rawconn.Close()
			return
		}

		conn := NewConn(rawconn)

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

		channel := communication.NewChannel(id, conn)
		channel.SetWriteWait(s.options.writeWait)
		channel.SetReadWait(s.options.readWait)
		s.Add(channel)

		go func(ch communication.Channel) {
			err := ch.ReadLoop(s.MessageListener)
			if err != nil {
				log.Info(err)
			}
			s.Remove(ch.ID())
			err = s.Disconnect(ch.ID())
			if err != nil {
				log.Warn(err)
			}
			ch.Close()
		}(channel)
	})

	log.Infoln("started")
	return http.ListenAndServe(s.listen, mux)
}

func (s *Server) SetAcceptor(acceptor communication.Acceptor) {
	s.Acceptor = acceptor
}

func (s *Server) SetMessageListener(messageListener communication.MessageListener) {
	s.MessageListener = messageListener
}

func (s *Server) SetStateListener(stateListener communication.StateListener) {
	s.StateListener = stateListener
}

func (s *Server) SetReadWait(duration time.Duration) {
	s.options.readWait = duration
}

func (s *Server) SetChannelMap(channelMap communication.ChannelMap) {
	s.ChannelMap = channelMap
}

func (s *Server) Push(channelId string, data []byte) error {
	ch, ok := s.ChannelMap.Get(channelId)
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

func NewServer(listen string, service communication.ServiceRegistration) communication.Server {
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
