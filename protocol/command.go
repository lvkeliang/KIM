package protocol

import "net"

type OpCode byte

// 直接搬来websocket的OpCode, 减少转换逻辑
const (
	OpContinuation OpCode = 0x0
	OpText         OpCode = 0x1
	OpBinary       OpCode = 0x2
	OpClose        OpCode = 0x8
	OpPing         OpCode = 0x9
	OpPong         OpCode = 0xa
)

// Frame 抽象websocket的Frame为一个接口, 用于封包和拆包
type Frame interface {
	SetOpCode(opcode OpCode)
	GetOpCode() OpCode
	SetPayload([]byte)
	GetPayload() []byte
}

// Conn 继承net.Conn, 进行二次封装, 方便读和写
type Conn interface {
	net.Conn
	// ReadFrame websocket/tcp两种协议的拆包
	ReadFrame() (Frame, error)
	// WriteFrame websocket/tcp两种协议的封包
	WriteFrame(OpCode, []byte) error
	Flush() error
}
