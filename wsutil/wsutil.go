package wsutil

import (
	"github.com/gobwas/pool/pbytes"
	"github.com/gobwas/ws"
	"io"
)

func writeFrame(w io.Writer, s ws.State, op ws.OpCode, fin bool, p []byte) error {
	var frame ws.Frame
	if s.ClientSide() { //区分客户端逻辑
		// Should copy bytes to prevent corruption of caller data.
		payload := pbytes.GetLen(len(p))
		defer pbytes.Put(payload)

		copy(payload, p)

		frame = ws.NewFrame(op, fin, payload)
		frame = ws.MaskFrameInPlace(frame) //对payload做了处理
	} else {
		frame = ws.NewFrame(op, fin, p)
	}

	return ws.WriteFrame(w, frame)
}
