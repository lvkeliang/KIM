package login

import (
	"KIM/communication"
	"KIM/communication/websocket"
	"KIM/logger"
	"KIM/protocol"
	"KIM/protocol/protoImpl"
	"KIM/token"
	"bytes"
	"context"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

type ClientDialer struct {
	AppSecret string
}

func (d *ClientDialer) DialAndHandshake(ctx communication.DialerContext) (net.Conn, error) {
	// 1. 拨号
	conn, _, _, err := ws.Dial(context.TODO(), ctx.Address)
	if err != nil {
		return nil, err
	}
	if d.AppSecret == "" {
		d.AppSecret = token.DefaultSecret
	}
	// 2. 直接使用封装的JWT包生成一个token
	tk, err := token.Generate(d.AppSecret, &token.Token{
		Account: ctx.Id,
		App:     "kim",
		Exp:     time.Now().AddDate(0, 0, 1).Unix(),
	})
	if err != nil {
		return nil, err
	}
	// 3. 发送一条CommandLoginSignIn消息
	loginreq := protocol.NewLogicPkt(protocol.CommandLoginSignIn).WriteBody(&protoImpl.LoginReq{
		Token: tk,
	})
	err = wsutil.WriteClientBinary(conn, protocol.MarshalPacket(loginreq))
	if err != nil {
		return nil, err
	}

	// wait resp
	_ = conn.SetReadDeadline(time.Now().Add(ctx.Timeout))
	frame, err := ws.ReadFrame(conn)
	if err != nil {
		return nil, err
	}
	ack, err := protocol.MustUnMarshalLogicPkt(bytes.NewBuffer(frame.Payload))
	if err != nil {
		return nil, err
	}
	// 4. 判断是否登录成功
	if ack.Status != protoImpl.Status_Success {
		return nil, fmt.Errorf("login failed: %v", &ack.Header)
	}
	var resp = new(protoImpl.LoginResp)
	_ = ack.ReadBody(resp)

	logger.Debug("logined ", resp.GetChannelId())
	return conn, nil
}

func Login(wsurl, account string, appSecrets ...string) (communication.Client, error) {
	cli := websocket.NewClient(account, "unittest", websocket.ClientOptions{})
	secret := token.DefaultSecret
	if len(appSecrets) > 0 {
		secret = appSecrets[0]
	}
	cli.SetDialer(&ClientDialer{
		AppSecret: secret,
	})
	err := cli.Connect(wsurl)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// const wsurl = "ws://119.3.4.216:8000"
const wsurl = "ws://localhost:8000"

func Test_login(t *testing.T) {
	cli, err := Login(wsurl, "test1")
	assert.Nil(t, err)
	time.Sleep(time.Second * 2)
	cli.Close()
}
