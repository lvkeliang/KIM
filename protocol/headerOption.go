package protocol

import "KIM/protocol/protoImpl"

// HeaderOption LogicPacket的Header初始化选项
type HeaderOption func(*protoImpl.Header)

func WithStatus(status protoImpl.Status) HeaderOption {
	return func(h *protoImpl.Header) {
		h.Status = status
	}
}

func WithSeq(seq uint32) HeaderOption {
	return func(h *protoImpl.Header) {
		h.Sequence = seq
	}
}

// WithChannel set channelID
func WithChannel(channelID string) HeaderOption {
	return func(h *protoImpl.Header) {
		h.ChannelId = channelID
	}
}

// WithDest WithDest
func WithDest(dest string) HeaderOption {
	return func(h *protoImpl.Header) {
		h.Dest = dest
	}
}
