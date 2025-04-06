package serv

import (
	"KIM/communication"
	"KIM/container"
	"KIM/logger"
	"KIM/protocol"
	"KIM/protocol/protoImpl"
	"KIM/token"
	"KIM/util"
	"bytes"
	"errors"
	"fmt"
	"github.com/prometheus/common/log"
	"regexp"
	"time"
)

type Handler struct {
	ServiceId string
	AppSecret string
}

// Disconnect 登出、断开连接、心跳超时会发出连接断开通知
func (h *Handler) Disconnect(id string) error {
	log.Infof("disconnect %s", id)

	logout := protocol.NewLogicPkt(protocol.CommandLoginSignOut, protocol.WithChannel(id))
	err := container.Forward(protocol.SNLogin, logout)
	if err != nil {
		logger.WithFields(logger.Fields{
			"module": "handler",
			"id":     id,
		}).Error(err)
	}
	return nil
}

// Receive 网关接收用户发送过来的消息
func (h *Handler) Receive(ag communication.Agent, payload []byte) {
	buf := bytes.NewBuffer(payload)
	packet, err := protocol.UnMarshalPacket(buf)
	if err != nil {
		return
	}

	// 收到basicPkt心跳，直接返回Pong
	if basicPkt, ok := packet.(*protocol.BasicPkt); ok {
		if basicPkt.Code == protocol.CodePing {
			_ = ag.Push(protocol.MarshalPacket(&protocol.BasicPkt{Code: protocol.CodePong}))
		}
		return
	}

	// 收到logicPkt，转发给逻辑服务处理
	if logicPkt, ok := packet.(*protocol.LogicPkt); ok {
		logicPkt.ChannelId = ag.ID()

		err = container.Forward(logicPkt.ServiceName(), logicPkt)
		if err != nil {
			logger.WithFields(logger.Fields{
				"module": "handler",
				"id":     ag.ID(),
				"cmd":    logicPkt.Command,
				"dest":   logicPkt.Dest,
			}).Error(err)
		}
	}
}

func (h *Handler) Accept(conn protocol.Conn, timeout time.Duration) (string, error) {
	log := logger.WithFields(logger.Fields{
		"ServiceID": h.ServiceId,
		"module":    "Handler",
		"handler":   "Accept",
	})
	log.Infoln("enter")

	// 读取登录包
	_ = conn.SetReadDeadline(time.Now().Add(timeout))
	frame, err := conn.ReadFrame()
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(frame.GetPayload())
	req, err := protocol.MustUnMarshalLogicPkt(buf)
	if err != nil {
		return "", err
	}

	// 必须是登录包
	if req.Command != protocol.CommandLoginSignIn {
		resp := protocol.NewLogicPktFromHeader(&req.Header)
		resp.Status = protoImpl.Status_InvalidCommand
		_ = conn.WriteFrame(protocol.OpBinary, protocol.MarshalPacket(resp))
		return "", errors.New("must be a InvalidCommand command")
	}

	// 反序列化Body
	var login protoImpl.LoginReq
	err = req.ReadBody(&login)
	if err != nil {
		return "", err
	}

	// 使用默认的DefaultSecret解析token
	tk, err := token.Parse(token.DefaultSecret, login.Token)
	if err != nil {
		// 如果token无效，就返回一个Unauthorized消息
		resp := protocol.NewLogicPktFromHeader(&req.Header)
		resp.Status = protoImpl.Status_Unauthorized
		_ = conn.WriteFrame(protocol.OpBinary, protocol.MarshalPacket(resp))
		return "", err
	}

	// 生成一个全局唯一的ChannelID
	id := generateChannelID(h.ServiceId, tk.Account)

	req.ChannelId = id
	req.WriteBody(&protoImpl.Session{
		Account:   tk.Account,
		ChannelId: id,
		GateId:    h.ServiceId,
		App:       tk.App,
		RemoteIP:  getIP(conn.RemoteAddr().String()),
	})

	// 把login转发给Login服务
	err = container.Forward(protocol.SNLogin, req)
	if err != nil {
		return "", err
	}
	return id, err
}

func generateChannelID(serviceID, account string) string {
	return fmt.Sprintf("%s_%s_%d", serviceID, account, util.Seq.Next())
}

var ipExp = regexp.MustCompile(string("\\:[0-9]+$"))

func getIP(remoteAddr string) string {
	if remoteAddr == "" {
		return ""
	}
	return ipExp.ReplaceAllString(remoteAddr, "")
}
