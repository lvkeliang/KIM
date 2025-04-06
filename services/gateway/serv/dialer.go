package serv

import (
	"KIM/communication"
	"KIM/communication/tcp"
	"KIM/logger"
	"KIM/protocol"
	"KIM/protocol/protoImpl"
	"github.com/golang/protobuf/proto"
	"net"
	"strings"
)

// TcpDialer 网关Dialer，用于服务间的连接建立
type TcpDialer struct {
	ServiceId string
}

func (d *TcpDialer) DialAndHandshake(ctx communication.DialerContext) (net.Conn, error) {
	// 1 调用net.Dial拨号
	address := ctx.Address
	if strings.Contains(address, "://") {
		// 去除协议前缀
		address = strings.Split(address, "://")[1]
	}
	conn, err := net.DialTimeout("tcp", address, ctx.Timeout)
	if err != nil {
		return nil, err
	}

	req := &protoImpl.InnerHandshakeReq{
		ServiceId: d.ServiceId,
	}
	logger.Infof("send handshakeReq %v", req)
	bts, _ := proto.Marshal(req)

	// 2. 发送ServiceId
	err = tcp.WriteFrame(conn, protocol.OpBinary, bts)
	if err != nil {
		return nil, err
	}
	// 3. return conn
	return conn, nil
}

func NewDialer(serviceid string) communication.Dialer {
	return &TcpDialer{
		ServiceId: serviceid,
	}
}
