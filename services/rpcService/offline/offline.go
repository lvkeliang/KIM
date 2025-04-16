package offline

import (
	"KIM/config"
	"KIM/database"
	"KIM/protocol/rpc"
	"KIM/services/rpcService/util"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type ServiceHandlerOffline struct {
	BaseDb    *gorm.DB
	MessageDb *gorm.DB
	Cache     *redis.Client
	Idgen     *database.IDGenerator
}

// InsertUserMessage 离线服务，向离线队列添加一条单聊消息
func (h *ServiceHandlerOffline) InsertUserMessage(c *gin.Context) {
	var req rpc.InsertMessageReq
	if err := util.AutoBind(c, &req); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	messageId, err := h.insertUserMessage(&req)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	util.Negotiate(c, http.StatusOK, &rpc.InsertMessageResp{
		MessageId: messageId,
	})
}

func (h *ServiceHandlerOffline) insertUserMessage(req *rpc.InsertMessageReq) (int64, error) {
	messageId := h.Idgen.Next().Int64()

	// 扩散写
	// direction,0表示accountB是发送方，1表示accountA是发送方
	idxs := make([]database.MessageIndex, 2)
	idxs[0] = database.MessageIndex{
		ID:        h.Idgen.Next().Int64(),
		AccountA:  req.Dest,
		AccountB:  req.Sender,
		Direction: 0,
		MessageID: messageId,
		SendTime:  req.SendTime,
	}
	idxs[1] = database.MessageIndex{
		ID:        h.Idgen.Next().Int64(),
		AccountA:  req.Sender,
		AccountB:  req.Dest,
		Direction: 1,
		MessageID: messageId,
		SendTime:  req.SendTime,
	}

	messageContent := database.MessageContent{
		ID:       messageId,
		Type:     byte(req.Message.Type),
		Body:     req.Message.Body,
		Extra:    req.Message.Extra,
		SendTime: req.SendTime,
	}

	err := h.MessageDb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&messageContent).Error; err != nil {
			return err
		}
		if err := tx.Create(&idxs).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return messageId, nil
}

// InsertGroupMessage 离线服务，向离线队列添加一条群聊消息
func (h *ServiceHandlerOffline) InsertGroupMessage(c *gin.Context) {
	var req rpc.InsertMessageReq
	if err := util.AutoBind(c, &req); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	messageId, err := h.insertGroupMessage(&req)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	util.Negotiate(c, http.StatusOK, &rpc.InsertMessageResp{
		MessageId: messageId,
	})
}

func (h *ServiceHandlerOffline) insertGroupMessage(req *rpc.InsertMessageReq) (int64, error) {
	messageId := h.Idgen.Next().Int64()

	var members []database.GroupMember
	err := h.BaseDb.Where(&database.GroupMember{Group: req.Dest}).Find(&members).Error
	if err != nil {
		return 0, err
	}

	// 扩散写
	// direction,0表示accountB是发送方，1表示accountA是发送方
	idxs := make([]database.MessageIndex, len(members))
	for i, m := range members {
		idxs[i] = database.MessageIndex{
			ID:        h.Idgen.Next().Int64(),
			AccountA:  m.Account,
			AccountB:  req.Sender,
			Direction: 0,
			MessageID: messageId,
			Group:     m.Group,
			SendTime:  req.SendTime,
		}
		if m.Account == req.Sender {
			idxs[i].Direction = 1
		}
	}

	messageContent := database.MessageContent{
		ID:       messageId,
		Type:     byte(req.Message.Type),
		Body:     req.Message.Body,
		Extra:    req.Message.Extra,
		SendTime: req.SendTime,
	}

	err = h.MessageDb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&messageContent).Error; err != nil {
			return err
		}
		if err := tx.Create(&idxs).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return messageId, nil
}

// MessageAck 离线服务，通过ACK确认接收到消息，只用确认接受到的最后一条消息，表示该消息及其之前的消息都已经确认
func (h *ServiceHandlerOffline) MessageAck(c *gin.Context) {
	var req rpc.AckMessageReq
	if err := util.AutoBind(c, &req); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	// 保存到redis
	err := setMessageAck(h.Cache, req.Account, req.MessageId)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
}

func setMessageAck(cache *redis.Client, account string, messageId int64) error {
	if messageId == 0 {
		return nil
	}
	key := database.KeyMessageAckIndex(account)
	return cache.Set(key, messageId, config.OfflineReadIndexExpiresIn).Err()
}

// GetOfflineMessageIndex 获取消息索引
func (h *ServiceHandlerOffline) GetOfflineMessageIndex(c *gin.Context) {
	var req rpc.AckMessageReq
	if err := util.AutoBind(c, &req); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	msgId := req.MessageId
	start, err := h.getSentTime(req.Account, req.MessageId)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	var indexes []*rpc.MessageIndex
	tx := h.MessageDb.Model(&database.MessageIndex{}).Select("send_time", "account_b", "direction", "message_id", "group")
	//err = tx.Where("account_a=? and send_time>? and direction=?", req.Account, start, 0).Order("send_time asc").Limit(config.OfflineSyncIndexCount).Find(&indexes).Error
	err = tx.Where("account_a=? and send_time>?", req.Account, start).Order("send_time asc").Limit(config.OfflineSyncIndexCount).Find(&indexes).Error
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	err = setMessageAck(h.Cache, req.Account, msgId)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	util.Negotiate(c, http.StatusOK, &rpc.GetOfflineMessageIndexResp{
		List: indexes,
	})
}

func (h *ServiceHandlerOffline) getSentTime(account string, messageId int64) (int64, error) {
	//冷启动情况
	if messageId == 0 {
		key := database.KeyMessageAckIndex(account)
		messageId, _ = h.Cache.Get(key).Int64() //如果一次都没有发ACK包，这里就是0
	}
	var start int64
	if messageId > 0 {
		//根据消息ID读取此条消息的发送时间
		var content database.MessageContent
		err := h.MessageDb.Select("send_time").First(&content, messageId).Error
		if err != nil {
			// 如果此条消息不存在，返回最近一天
			start = time.Now().AddDate(0, 0, -1).UnixNano()
		} else {
			start = content.SendTime
		}
	}
	//返回默认的离线消息过期时间
	earliestKeepTime := time.Now().AddDate(0, 0, -1*config.OfflineMessageExpiresIn).UnixNano()
	if start == 0 || start < earliestKeepTime {
		start = earliestKeepTime
	}
	return start, nil
}

// GetOfflineMessageContent 下载消息内容
func (h *ServiceHandlerOffline) GetOfflineMessageContent(c *gin.Context) {
	var req rpc.GetOfflineMessageContentReq
	if err := util.AutoBind(c, &req); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	mlen := len(req.MessageIds)
	if mlen > config.MessageMaxCountPerPage {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "too many MessageIds"})
		return
	}

	var contents []*rpc.Message
	err := h.MessageDb.Model(&database.MessageContent{}).Where(req.MessageIds).Find(&contents).Error
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	util.Negotiate(c, http.StatusOK, &rpc.GetOfflineMessageContentResp{
		List: contents,
	})
}
