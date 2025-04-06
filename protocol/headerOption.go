package protocol

import "KIM/protocol/protoImpl"

// HeaderOption HeaderOption
type HeaderOption func(*protoImpl.Header)

// WithStatus WithStatus
func WithStatus(status protoImpl.Status) HeaderOption {
	return func(h *protoImpl.Header) {
		h.Status = status
	}
}

// WithSeq WithSeq
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
