package login

import (
	"KIM/logger"
	"KIM/protocol/protoImpl"
	"KIM/services"
)

type LoginHandler struct {
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (h *LoginHandler) DoSysLogin(ctx services.RouterContext) {
	log := logger.WithField("func", "DoSysLogin")
	// 序列化
	var session protoImpl.Session
	if err := ctx.ReadBody(&session); err != nil {
		_ = ctx.RespWithError(protoImpl.Status_InvalidPacketBody, err)
		return
	}

	log.Infof("do login of %v ", session.String())
	// 检查当前账号是否已经登录在其它地方
	old, err := ctx.GetLocation(session.Account, "")
	if err != nil && err != services.ErrSessionNil {
		_ = ctx.RespWithError(protoImpl.Status_SystemException, err)
		return
	}

	if old != nil {
		// 通知该连接下线
		_ = ctx.Dispatch(&protoImpl.KickOutNotify{
			ChannelId: old.ChannelId,
		}, old)
	}

	// 添加到会话管理
	err = ctx.Add(&session)
	if err != nil {
		_ = ctx.RespWithError(protoImpl.Status_SystemException, err)
		return
	}

	// 返回一个登录成功消息
	var resp = &protoImpl.LoginResp{
		ChannelId: session.ChannelId,
		Account:   session.Account,
	}

	_ = ctx.Resp(protoImpl.Status_Success, resp)
}

// DoSysLogout 登出靠网关handler的Disconnect发出连接断开通知
func (h *LoginHandler) DoSysLogout(ctx services.RouterContext) {
	logger.WithField("func", "DoSysLogout").Infof("do Logout of %s %s", ctx.Session().GetChannelId(), ctx.Session().GetAccount())

	err := ctx.Delete(ctx.Session().GetAccount(), ctx.Session().GetChannelId())
	if err != nil {
		_ = ctx.RespWithError(protoImpl.Status_SystemException, err)
		return
	}

	_ = ctx.Resp(protoImpl.Status_Success, nil)
}
