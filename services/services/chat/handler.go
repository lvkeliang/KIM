package chat

import (
	"KIM/protocol/protoImpl"
	"KIM/protocol/rpc"
	"KIM/services"
	"KIM/services/services/rpcServiceClient/group"
	"KIM/services/services/rpcServiceClient/offline"
	"errors"
	"time"
)

type ChatHandler struct {
	offlineService offline.Offline
	groupService   group.Group
}

func NewChatHandler(offlineService offline.Offline, groupService group.Group) *ChatHandler {
	return &ChatHandler{offlineService: offlineService, groupService: groupService}
}

var ErrNoDestination = errors.New("dest is empty")

func (h *ChatHandler) DoUserTalk(ctx services.RouterContext) {
	if ctx.Header().Dest == "" {
		_ = ctx.RespWithError(protoImpl.Status_NoDestination, ErrNoDestination)
		return
	}

	var req protoImpl.MessageReq
	if err := ctx.ReadBody(&req); err != nil {
		_ = ctx.RespWithError(protoImpl.Status_InvalidPacketBody, err)
		return
	}

	//接收方寻址
	receiver := ctx.Header().GetDest()
	loc, err := ctx.GetLocation(receiver, "")
	if err != nil && err != services.ErrSessionNil {
		_ = ctx.RespWithError(protoImpl.Status_SystemException, err)
	}

	sendTime := time.Now().UnixNano() //TODO:此处使用的是服务器本地时钟，在IM系统中可能有问题，需要换成全局时钟

	//保存离线消息
	resp, err := h.offlineService.InsertUser(ctx.Session().GetApp(), &rpc.InsertMessageReq{
		Sender:   ctx.Session().GetAccount(),
		Dest:     receiver,
		SendTime: sendTime,
		Message: &rpc.Message{
			Type:  req.GetType(),
			Body:  req.GetBody(),
			Extra: req.GetExtra(),
		},
	})

	if err != nil {
		_ = ctx.RespWithError(protoImpl.Status_SystemException, err)
		return
	}

	// 如果接收方在线，推送一条消息过去
	if loc != nil {
		if err = ctx.Dispatch(&protoImpl.MessagePush{
			MessageId: resp.MessageId,
			Type:      req.GetType(),
			Body:      req.GetBody(),
			Extra:     req.GetExtra(),
			Sender:    ctx.Session().GetAccount(),
			SendTime:  sendTime,
		}, loc); err != nil {
			_ = ctx.RespWithError(protoImpl.Status_SystemException, err)
			return
		}
	}

	//ACK: resp一条Status_Success消息
	_ = ctx.Resp(protoImpl.Status_Success, &protoImpl.MessageResp{
		MessageId: resp.MessageId,
		SendTime:  sendTime,
	})
}

func (h *ChatHandler) DoGroupTalk(ctx services.RouterContext) {
	if ctx.Header().Dest == "" {
		_ = ctx.RespWithError(protoImpl.Status_NoDestination, ErrNoDestination)
		return
	}

	var req protoImpl.MessageReq
	if err := ctx.ReadBody(&req); err != nil {
		_ = ctx.RespWithError(protoImpl.Status_InvalidPacketBody, err)
		return
	}

	//接收群聊寻址, 群聊的Dest不再是user account, 而是群ID
	group := ctx.Header().GetDest()

	sendTime := time.Now().UnixNano() //TODO:此处使用的是服务器本地时钟，在IM系统中可能有问题，需要换成全局时钟，可以使用raft算法的ETCD这类kv服务来决定消息提交顺序4

	//保存离线消息
	resp, err := h.offlineService.InsertGroup(ctx.Session().GetApp(), &rpc.InsertMessageReq{
		Sender:   ctx.Session().GetAccount(),
		Dest:     group,
		SendTime: sendTime,
		Message: &rpc.Message{
			Type:  req.GetType(),
			Body:  req.GetBody(),
			Extra: req.GetExtra(),
		},
	})
	if err != nil {
		_ = ctx.RespWithError(protoImpl.Status_SystemException, err)
		return
	}

	// 读取成员列表
	membersResp, err := h.groupService.Members(ctx.Session().GetApp(), &rpc.GroupMembersReq{
		GroupId: group,
	})
	if err != nil {
		_ = ctx.RespWithError(protoImpl.Status_SystemException, err)
		return
	}
	var members = make([]string, len(membersResp.Users))
	for i, user := range membersResp.Users {
		members[i] = user.Account
	}
	// 批量寻址
	locs, err := ctx.GetLocations(members...)
	if err != nil && err != services.ErrSessionNil {
		_ = ctx.RespWithError(protoImpl.Status_SystemException, err)
		return
	}

	// 批量推送给群成员
	if len(locs) > 0 {
		if err = ctx.Dispatch(&protoImpl.MessagePush{
			MessageId: resp.MessageId,
			Type:      req.GetType(),
			Body:      req.GetBody(),
			Extra:     req.GetExtra(),
			Sender:    ctx.Session().GetAccount(),
			SendTime:  sendTime,
		}, locs...); err != nil {
			_ = ctx.RespWithError(protoImpl.Status_SystemException, err)
			return
		}
	}

	//ACK: resp一条Status_Success消息
	_ = ctx.Resp(protoImpl.Status_Success, &protoImpl.MessageResp{
		MessageId: resp.MessageId,
		SendTime:  sendTime,
	})
}

func (h *ChatHandler) DoTalkAck(ctx services.RouterContext) {
	var req protoImpl.MessageAckReq
	if err := ctx.ReadBody(&req); err != nil {
		_ = ctx.RespWithError(protoImpl.Status_InvalidPacketBody, err)
		return
	}
	err := h.offlineService.SetAck(ctx.Session().GetApp(), &rpc.AckMessageReq{
		Account:   ctx.Session().GetAccount(),
		MessageId: req.GetMessageId(),
	})
	if err != nil {
		_ = ctx.RespWithError(protoImpl.Status_SystemException, err)
		return
	}
	_ = ctx.Resp(protoImpl.Status_Success, nil)
}
