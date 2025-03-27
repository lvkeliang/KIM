package main

import (
	"KIM/kimwebsocketdemo/protocol"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"sync"
	"time"
)

const heartBeatDeadline = time.Minute * 2

func main() {
	RunServerStart(context.Background(), &ServerStartOptions{
		id:     "Server1",
		listen: ":8080",
	}, "v0.1")
}

type ServerStartOptions struct {
	id     string
	listen string
}

type Server struct {
	once    sync.Once
	id      string
	address string
	sync.RWMutex
	//会话列表
	users map[string]net.Conn
}

func RunServerStart(ctx context.Context, opts *ServerStartOptions, version string) error {
	server := NewServer(opts.id, opts.listen)
	defer server.Shutdown()
	return server.Start()
}

func NewServer(id, address string) *Server {
	return newServer(id, address)
}

func newServer(id, address string) *Server {
	return &Server{
		id:      id,
		address: address,
		users:   make(map[string]net.Conn, 100),
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	log := logrus.WithFields(logrus.Fields{
		"module": "Server",
		"listen": s.address,
		"id":     s.id,
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 升级
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			conn.Close()
			return
		}

		// 读取userId
		user := r.URL.Query().Get("user")
		if user == "" {
			conn.Close()
			return
		}

		// 添加到会话管理
		old, ok := s.addUser(user, conn)
		if ok {
			// 断开旧连接
			old.Close()
		}
		log.Infof("user %s in", user)

		go func(user string, conn net.Conn) {
			// 读取消息
			err := s.readloop(user, conn)
			if err != nil {
				log.Error(err)
			}
			conn.Close()

			// 连接断开，删除用户
			s.delUser(user)
			log.Infof("connection of %s closed", user)
		}(user, conn)
	})

	log.Infoln("started")
	return http.ListenAndServe(s.address, mux)
}

func (s *Server) addUser(user string, conn net.Conn) (net.Conn, bool) {
	s.Lock()
	defer s.Unlock()
	old, ok := s.users[user] //返回旧连接
	s.users[user] = conn     //缓存新连接
	return old, ok
}

func (s *Server) delUser(user string) {
	s.Lock()
	defer s.Unlock()
	delete(s.users, user)
}

// Shutdown 关闭Server
func (s *Server) Shutdown() {
	s.once.Do(func() {
		s.Lock()
		defer s.Unlock()
		for _, conn := range s.users {
			conn.Close()
		}
	})
}

func (s *Server) readloop(user string, conn net.Conn) error {
	for {
		// 要求客户端必须在指定时间内发送一条消息过来，可以是ping或其他数据包
		_ = conn.SetReadDeadline(time.Now().Add(heartBeatDeadline))

		// 读取一帧信息
		frame, err := ws.ReadFrame(conn)
		if err != nil {
			return err
		}

		if frame.Header.OpCode == ws.OpPing {
			// 返回一个Pone消息
			_ = wsutil.WriteServerMessage(conn, ws.OpPong, nil)
			continue
		}

		if frame.Header.OpCode == ws.OpClose {
			return errors.New("remote side close the conn")
		}

		if frame.Header.Masked {
			ws.Cipher(frame.Payload, frame.Header.Mask, 0)
		}

		// 读取文本帧内容
		if frame.Header.OpCode == ws.OpText {
			go s.handle(user, string(frame.Payload))
		} else if frame.Header.OpCode == ws.OpBinary {
			go s.handleBinary(user, frame.Payload)
		}
	}
}

func (s *Server) handle(user string, message string) {
	logrus.Infof("recv message %s from %s", message, user)
	s.Lock()
	defer s.Unlock()
	// s.users不线程安全, 需加读锁
	// ws.WriteFrame(conn, f)不线程安全，需加写锁
	// 所以这里简单的对整个handle加了一个锁
	broadcast := fmt.Sprintf("%s -- FROM %s", message, user)
	for u, conn := range s.users {
		if u == user {
			continue
		}
		logrus.Infof("send to %s : %s", u, broadcast)
		err := s.writeText(conn, broadcast)
		if err != nil {
			logrus.Errorf("write to %s failed, error: %v", user, err)
		}
	}
}

func (s *Server) writeText(conn net.Conn, message string) error {
	// 创建文本数据
	f := ws.NewTextFrame([]byte(message))
	return ws.WriteFrame(conn, f)
}

func (s *Server) handleBinary(user string, message []byte) {
	logrus.Infof("recv message %v from %s", message, user)
	s.RLock()
	defer s.RUnlock()

	// handle Ping请求
	i := 0
	command := binary.BigEndian.Uint32(message[i : i+2])
	i += 2
	payloadLen := binary.BigEndian.Uint32(message[i : i+4])
	logrus.Infof("command: %v payloadLen: %v", command, payloadLen)
	if command == protocol.CommandPing {
		u := s.users[user]
		// 响应 pong
		err := wsutil.WriteServerBinary(u, []byte{0, protocol.CommandPong, 0, 0, 0, 0})
		if err != nil {
			logrus.Errorf("write to %s failed, error: %v", user, err)
		}

	}
}
