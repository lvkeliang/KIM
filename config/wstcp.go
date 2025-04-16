package config

import "time"

// 通信层ws和tcp的等待时延
const (
	DefaultLoginWait = time.Second * 10
	DefaultReadWait  = time.Minute * 3
	DefaultWriteWait = time.Second * 10
	DefaultHeartbeat = time.Second * 55
)

const (
	OfflineReadIndexExpiresIn = time.Hour * 24 * 30 // 读索引在缓存中的过期时间
	OfflineSyncIndexCount     = 2000                //单次同步消息索引的数量
	OfflineMessageExpiresIn   = 15                  // 离线消息过期时间
	MessageMaxCountPerPage    = 200                 // 同步消息内容时每页的最大数据
)
