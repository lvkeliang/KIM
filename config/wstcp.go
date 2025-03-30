package config

import "time"

// 通信层ws和tcp的等待时延
const (
	DefaultLoginWait = time.Second * 10
	DefaultReadWait  = time.Minute * 3
	DefaultWriteWait = time.Second * 10
	DefaultHeartbeat = time.Second * 55
)
