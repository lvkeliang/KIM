package config

import "time"

const (
	DefaultLoginWait = time.Second * 10
	DefaultReadWait  = time.Minute * 3
	DefaultWriteWait = time.Second * 10
	DefaultHeartbeat = time.Second * 55
)
