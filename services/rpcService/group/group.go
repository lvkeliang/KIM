package group

import (
	"KIM/database"
	"KIM/protocol/rpc"
	"KIM/services/rpcService/util"
	"errors"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
	"net/http"
)

type ServiceHandlerGroup struct {
	BaseDb    *gorm.DB
	MessageDb *gorm.DB
	Cache     *redis.Client
	Idgen     *database.IDGenerator
}

// GroupGet 获取群聊信息
func (h *ServiceHandlerGroup) GroupGet(c *gin.Context) {
	groupId := c.Param("id")
	if groupId == "" {
		util.ErrorResponse(c, http.StatusBadRequest, errors.New("group is null"))
		return
	}
	id, err := h.Idgen.ParseBase36(groupId)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, errors.New("group is invalid: "+groupId))
		return
	}

	var group database.Group
	err = h.BaseDb.First(&group, id.Int64()).Error
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, errors.New("group is invalid: "+groupId))
		return
	}

	util.Negotiate(c, http.StatusOK, &rpc.GetGroupResp{
		Id:           groupId,
		Name:         group.Name,
		Avatar:       group.Avatar,
		Introduction: group.Introduction,
		Owner:        group.Owner,
		CreatedAt:    group.CreatedAt.Unix(),
	})
}

// GroupCreate 创建群聊
func (h *ServiceHandlerGroup) GroupCreate(c *gin.Context) {
	app := c.Param("app")
	var req rpc.CreateGroupReq
	if err := util.AutoBind(c, &req); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	req.App = app
	groupId, err := h.groupCreate(&req)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	util.Negotiate(c, http.StatusOK, &rpc.CreateGroupResp{GroupId: groupId.Base36()})
}

func (h *ServiceHandlerGroup) groupCreate(req *rpc.CreateGroupReq) (snowflake.ID, error) {
	groupId := h.Idgen.Next()
	g := &database.Group{
		Model:        database.Model{ID: groupId.Int64()},
		Group:        groupId.Base36(),
		App:          req.App,
		Name:         req.Name,
		Owner:        req.Owner,
		Avatar:       req.Avatar,
		Introduction: req.Introduction,
	}

	members := make([]database.GroupMember, len(req.Members))
	for i, user := range req.Members {
		members[i] = database.GroupMember{
			Model:   database.Model{ID: h.Idgen.Next().Int64()},
			Account: user,
			Group:   groupId.Base36(),
		}
	}

	err := h.BaseDb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(g).Error; err != nil {
			return err
		}
		if err := tx.Create(&members).Error; err != nil {
			return err
		}
		// 返回nil提交事务
		return nil
	})
	if err != nil {
		return 0, err
	}
	return groupId, nil
}

// GroupJoin 加入群聊
func (h *ServiceHandlerGroup) GroupJoin(c *gin.Context) {
	var req rpc.JoinGroupReq
	if err := util.AutoBind(c, &req); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	gm := &database.GroupMember{
		Model:   database.Model{ID: h.Idgen.Next().Int64()},
		Account: req.Account,
		Group:   req.GroupId,
	}

	err := h.BaseDb.Create(gm).Error
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
}

// GroupQuit 退出群聊
func (h *ServiceHandlerGroup) GroupQuit(c *gin.Context) {
	var req rpc.QuitGroupReq
	if err := util.AutoBind(c, &req); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	gm := &database.GroupMember{
		Account: req.Account,
		Group:   req.GroupId,
	}

	err := h.BaseDb.Delete(&database.GroupMember{}, gm).Error
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	//c.JSON(http.StatusOK, "")
}

// GroupMembers 获取群成员
func (h *ServiceHandlerGroup) GroupMembers(c *gin.Context) {
	groupId := c.Param("id")
	if groupId == "" {
		util.ErrorResponse(c, http.StatusBadRequest, errors.New("group is null"))
		return
	}

	var members []database.GroupMember
	err := h.BaseDb.Order("Updated_At asc").Find(&members, database.GroupMember{Group: groupId}).Error
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	var users = make([]*rpc.Member, len(members))
	for i, m := range members {
		users[i] = &rpc.Member{
			Account:  m.Account,
			Alias:    m.Alias,
			JoinTime: m.CreatedAt.Unix(),
		}
	}

	util.Negotiate(c, http.StatusOK, &rpc.GroupMembersResp{Users: users})
}
