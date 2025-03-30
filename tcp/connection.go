package tcp

import (
	"KIM/protocol"
	"KIM/util"
	"io"
	"net"
)

// Frame tcp Frame
type Frame struct {
	OpCode  protocol.OpCode
	Payload []byte
}

// SetOpCode SetOpCode
func (f *Frame) SetOpCode(code protocol.OpCode) {
	f.OpCode = code
}

// GetOpCode GetOpCode
func (f *Frame) GetOpCode() protocol.OpCode {
	return f.OpCode
}

// SetPayload SetPayload
func (f *Frame) SetPayload(payload []byte) {
	f.Payload = payload
}

// GetPayload GetPayload
func (f *Frame) GetPayload() []byte {
	return f.Payload
}

// TcpConn Conn
type TcpConn struct {
	net.Conn
}

// NewConn NewConn
func NewConn(conn net.Conn) *TcpConn {
	return &TcpConn{
		Conn: conn,
	}
}

// ReadFrame ReadFrame
func (c *TcpConn) ReadFrame() (protocol.Frame, error) {
	opcode, err := util.ReadUint8(c.Conn)
	if err != nil {
		return nil, err
	}
	payload, err := util.ReadBytes(c.Conn)
	if err != nil {
		return nil, err
	}
	return &Frame{
		OpCode:  protocol.OpCode(opcode),
		Payload: payload,
	}, nil
}

// WriteFrame WriteFrame
func (c *TcpConn) WriteFrame(code protocol.OpCode, payload []byte) error {
	return WriteFrame(c.Conn, code, payload)
}

// Flush Flush
func (c *TcpConn) Flush() error {
	return nil
}

// WriteFrame write a frame to w
func WriteFrame(w io.Writer, code protocol.OpCode, payload []byte) error {
	if err := util.WriteUint8(w, uint8(code)); err != nil {
		return err
	}
	if err := util.WriteBytes(w, payload); err != nil {
		return err
	}
	return nil
}
