package container

import (
	"KIM/inter"
	"KIM/protocol/protoImpl"
	"hash/crc32"
)

// HashSelector 哈希选择器
type HashSelector struct{}

func (s *HashSelector) Lookup(header *protoImpl.Header, srvs []inter.Service) string {
	length := len(srvs)
	code := HashCode(header.ChannelId)
	return srvs[code%length].ServiceID()
}

func HashCode(key string) int {
	hash32 := crc32.NewIEEE()
	hash32.Write([]byte(key))
	return int(hash32.Sum32())
}
