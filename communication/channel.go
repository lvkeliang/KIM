package communication

import (
	"KIM/config"
	"KIM/logger"
	"KIM/protocol"
	"errors"
	"sync"
	"time"
)

// Channel Server会将将连接包装为Channel来管理
type Channel interface {
	protocol.Conn
	Agent
	Close() error                       // 重写net.Conn中的Close方法
	ReadLoop(lst MessageListener) error // 阻塞的方法, 把消息的读取和心跳处理封装在一起
	SetWriteWait(time.Duration)
	SetReadWait(time.Duration)
}

// ChannelMap 连接管理器, 管理Channel
type ChannelMap interface {
	Add(channel Channel)
	Remove(id string)
	Get(id string) (Channel, bool)
	All() []Channel
}

type ChannelsImpl struct {
	channels *sync.Map
}

// Add addChannel
func (ch *ChannelsImpl) Add(channel Channel) {
	if channel.ID() == "" {
		logger.WithFields(logger.Fields{
			"module": "ChannelsImpl",
		}).Error("channel id is required")
		return
	}

	ch.channels.Store(channel.ID(), channel)
}

// Remove addChannel
func (ch *ChannelsImpl) Remove(id string) {
	ch.channels.Delete(id)
}

// Get Get
func (ch *ChannelsImpl) Get(id string) (Channel, bool) {
	if id == "" {
		logger.WithFields(logger.Fields{
			"module": "ChannelsImpl",
		}).Error("channel id is required")
		return nil, false
	}

	val, ok := ch.channels.Load(id)
	if !ok {
		return nil, false
	}
	return val.(Channel), true
}

// All return channels
func (ch *ChannelsImpl) All() []Channel {
	arr := make([]Channel, 0)
	ch.channels.Range(func(key, val interface{}) bool {
		arr = append(arr, val.(Channel))
		return true
	})
	return arr
}

func NewChannels(num int) ChannelMap {
	return &ChannelsImpl{
		channels: new(sync.Map),
	}
}

// ChannelImpl 实现server的Channel
type ChannelImpl struct {
	sync.Mutex
	id string
	protocol.Conn
	writechan chan []byte
	once      sync.Once
	writeWait time.Duration
	readWait  time.Duration
	closed    *Event
}

func NewChannel(id string, conn protocol.Conn) Channel {
	log := logger.WithFields(logger.Fields{
		"module": "tcp_channel",
		"id":     id,
	})

	ch := ChannelImpl{
		id:        id,
		Conn:      conn,
		writechan: make(chan []byte, 5),
		writeWait: config.DefaultWriteWait,
		readWait:  config.DefaultReadWait,
		closed:    NewEvent(),
	}

	go func() {
		err := ch.writeLoop()
		if err != nil {
			log.Info(err)
		}
	}()
	return &ch
}

func (ch *ChannelImpl) writeLoop() error {
	for {
		select {
		case payload := <-ch.writechan:
			err := ch.WriteFrame(protocol.OpBinary, payload)
			if err != nil {
				return err
			}
			// 批量写
			chanlen := len(ch.writechan)
			for i := 0; i < chanlen; i++ {
				payload = <-ch.writechan
				err = ch.WriteFrame(protocol.OpBinary, payload)
				if err != nil {
					return err
				}
			}
			err = ch.Conn.Flush()
			if err != nil {
				return err
			}
		case <-ch.closed.DoneChan:
			return nil

		}
	}
}

func (ch *ChannelImpl) Push(payload []byte) error {
	if ch.closed.HasFired() {
		return errors.New("channel has closed")
	}

	// 异步写
	ch.writechan <- payload
	return nil
}

// WriteFrame 重写Conn的WriteFrame方法, 添加重置写超时的逻辑
func (ch *ChannelImpl) WriteFrame(code protocol.OpCode, payload []byte) error {
	_ = ch.Conn.SetWriteDeadline(time.Now().Add(ch.writeWait))
	return ch.Conn.WriteFrame(code, payload)
}

func (ch *ChannelImpl) ReadLoop(lst MessageListener) error {
	ch.Lock()
	defer ch.Unlock()
	log := logger.WithFields(logger.Fields{
		"struct": "ChannelImpl",
		"func":   "ReadLoop",
		"id":     ch.id,
	})

	for {
		_ = ch.SetReadDeadline(time.Now().Add(ch.readWait))

		frame, err := ch.ReadFrame()
		if err != nil {
			return err
		}
		if frame.GetOpCode() == protocol.OpClose {
			return errors.New("remote side close the channel")
		}
		if frame.GetOpCode() == protocol.OpPing {
			log.Trace("recv aping: resp wite a pong")
			_ = ch.WriteFrame(protocol.OpPong, nil)
			continue
		}
		payload := frame.GetPayload()
		if len(payload) == 0 {
			continue
		}
		go lst.Receive(ch, payload)
	}
}

func (ch *ChannelImpl) ID() string {
	return ch.id
}

func (ch *ChannelImpl) SetWriteWait(writeWait time.Duration) {
	if writeWait == 0 {
		return
	}
	ch.writeWait = writeWait
}
func (ch *ChannelImpl) SetReadWait(readWait time.Duration) {
	if readWait == 0 {
		return
	}
	ch.readWait = readWait
}
func (ch *ChannelImpl) Close() error { // 重写net.Conn中的Close方法
	if !ch.closed.HasFired() {
		return errors.New("channel has started")
	}
	close(ch.writechan)
	return nil
}
