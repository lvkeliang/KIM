package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/sirupsen/logrus"
	"net"
	"net/url"
	"time"
)

const heartbeatDuration = 3 * time.Second

func main() {
	err := run(context.Background(), &StartOptions{address: "127.0.0.1:8080", user: "lisi"})
	if err != nil {
		logrus.Error(err)
	}
}

type StartOptions struct {
	address string
	user    string
}
type handler struct {
	conn  net.Conn
	close chan struct{}
	recv  chan []byte

	heartbeat time.Duration
}

func connect(addr string) (*handler, error) {
	_, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	conn, _, _, err := ws.Dial(context.Background(), addr)
	if err != nil {
		return nil, err
	}

	h := handler{
		conn:      conn,
		close:     make(chan struct{}, 1),
		recv:      make(chan []byte, 10),
		heartbeat: time.Duration(heartbeatDuration),
	}

	go func() {
		err := h.readloop(conn)
		if err != nil {
			logrus.Warn(err)
		}

		// 通知上层
		h.close <- struct{}{}
	}()

	return &h, nil
}

func (h *handler) readloop(conn net.Conn) error {
	logrus.Info("readloop started")

	// 要在heartbeatDuration * 3的时间内可以读到数据
	err := h.conn.SetReadDeadline(time.Now().Add(h.heartbeat * 3))
	if err != nil {
		return err
	}
	for {
		frame, err := ws.ReadFrame(conn)
		if err != nil {
			return err
		}

		if frame.Header.OpCode == ws.OpPong {
			//重置超时时间
			_ = h.conn.SetReadDeadline(time.Now().Add(h.heartbeat * 3))
		}

		if frame.Header.OpCode == ws.OpClose {
			return errors.New("remote side close the channel")
		}

		if frame.Header.OpCode == ws.OpText {
			h.recv <- frame.Payload
		}
	}
}

func run(ctx context.Context, opts *StartOptions) error {
	url := fmt.Sprintf("ws://%s?user=%s", opts.address, opts.user)
	logrus.Info("connect to ", url)

	h, err := connect(url)
	if err != nil {
		return err
	}

	go func() {
		err := h.heartbeatloop()
		if err != nil {
			logrus.Warn("heartbeat ping failed")
			return
		}
	}()

	go func() {

		for msg := range h.recv {
			logrus.Info("Receive message:", string(msg))
		}
	}()

	tk := time.NewTicker(time.Second * 6)
	for {
		select {
		case <-tk.C:
			err := h.sendText("hello我是李四")
			if err != nil {
				logrus.Error(err)
			}
		case <-h.close:
			logrus.Printf("connnection closed")
			return nil
		}
	}
}

func (h *handler) sendText(msg string) error {
	logrus.Info("send message: ", msg)
	return wsutil.WriteClientText(h.conn, []byte(msg))
}

func (h *handler) heartbeatloop() error {
	logrus.Info("heartbeatloop started")
	tick := time.NewTicker(h.heartbeat)
	for range tick.C {
		logrus.Info("heartbeat ping...")
		// 发送一个平的心跳包给服务端
		if err := wsutil.WriteClientMessage(h.conn, ws.OpPing, nil); err != nil {
			return err
		}
	}
	return nil
}
