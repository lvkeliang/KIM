package communication

import (
	"sync"
	"sync/atomic"
)

// Event 是一个线程安全的通知机制，用于广播关闭信号
type Event struct {
	fired    uint32        // 原子标记，1表示已触发
	DoneChan chan struct{} // 关闭时广播
	mu       sync.Mutex    // 保护 DoneChan 的初始化
}

func NewEvent() *Event {
	return &Event{fired: 0, DoneChan: make(chan struct{})}
}

func (e *Event) Done() {
	if atomic.SwapUint32(&e.fired, 1) == 1 {
		return // 已经触发过
	}
	e.DoneChan <- struct{}{}
}

// HasFired 检查事件是否已触发
func (e *Event) HasFired() bool {
	return atomic.LoadUint32(&e.fired) == 1
}
