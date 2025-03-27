package websocket

import (
	"KIM/protocol"
	"github.com/gobwas/ws"
	"net"
)

// Frame 实现protocol的Frame接口
type Frame struct {
	raw ws.Frame
}

func (f *Frame) SetOpCode(code protocol.OpCode) {
	f.raw.Header.OpCode = ws.OpCode(code)
}

func (f *Frame) GetOpCode() protocol.OpCode {
	return protocol.OpCode(f.raw.Header.OpCode)
}

// SetPayload 注意, 这里没有使用Mask来编码，所以客户端发送消息时不能直接使用websocket.Conn这个对象
func (f *Frame) SetPayload(payload []byte) {
	f.raw.Payload = payload
}
func (f *Frame) GetPayload() []byte {
	if f.raw.Header.Masked {
		ws.Cipher(f.raw.Payload, f.raw.Header.Mask, 0)
	}
	f.raw.Header.Masked = false
	return f.raw.Payload
}

// WsConn 实现protocol的Conn接口
type WsConn struct {
	net.Conn
}

func NewConn(conn net.Conn) *WsConn {
	return &WsConn{
		Conn: conn,
	}
}

func (c *WsConn) ReadFrame() (protocol.Frame, error) {
	f, err := ws.ReadFrame(c.Conn)
	if err != nil {
		return nil, err
	}

	return &Frame{raw: f}, nil
}

// WriteFrame websocket/tcp两种协议的封包
func (c *WsConn) WriteFrame(code protocol.OpCode, payload []byte) error {
	// 由于我们的数据包大小不会超过一个websocket协议单个帧最大值，因此这里的fin全部使用true，不会把包拆分成多个
	f := ws.NewFrame(ws.OpCode(code), true, payload)
	return ws.WriteFrame(c.Conn, f)
}
func (c *WsConn) Flush() error {
	// TODO
	return nil
}
