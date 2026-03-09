package sms

import (
	"adb-backup/internal/config"
	"adb-backup/internal/device"
	"adb-backup/internal/web/base"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ConversationPageResult struct {
	List        []Conversation `json:"list"`         // 当前页会话列表
	HasMore     bool           `json:"has_more"`     // 是否还有更多会话
	CurrentPage int            `json:"current_page"` // 当前页码
	PageSize    int            `json:"page_size"`    // 每页条数
}

type MessageResult struct {
	Messages []Message `json:"messages"` // 最新消息列表（已反转，最新在最后）
	HasMore  bool      `json:"has_more"` // 是否有更多老消息
}

type MessageSendReq struct {
	base.ContextReq
	DeviceId string `json:"device_id" binding:"required,deviceIdConnect"`
	Address  string `json:"address" binding:"required"`
	Body     string `json:"body" binding:"required"`
	SubId    int    `json:"sub_id" binding:"oneof=0 1"`
}

func GetConversationsApiHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		// 1. 获取设备ID参数
		deviceId := c.GetString(base.ContextDeviceIdKey)
		pageStr := r.URL.Query().Get("page")
		pageSizeStr := r.URL.Query().Get("page_size")
		address := r.URL.Query().Get("address")

		// 默认分页参数
		page := 1
		if pageStr != "" {
			page, _ = strconv.Atoi(pageStr)
			if page < 1 {
				page = 1
			}
		}

		pageSize := 20
		if pageSizeStr != "" {
			pageSize, _ = strconv.Atoi(pageSizeStr)
			if pageSize < 1 || pageSize > 50 {
				pageSize = 20
			}
		}
		// 查询总条数
		totalCount, err := getConversationTotalCount(deviceId, address)
		if err != nil {
			base.RespJsonInternalServerError(c, "查询会话总数失败："+err.Error())
			return
		}

		// 计算偏移量
		offset := (page - 1) * pageSize
		// 2. 查询会话列表
		conversations, err := getConversations(deviceId, offset, pageSize, address)
		if err != nil {
			base.RespJsonInternalServerError(c, "获取会话列表失败："+err.Error())
			return
		}
		// 判断是否还有更多
		hasMore := (offset + pageSize) < int(totalCount)
		// 组装结果
		pageResult := ConversationPageResult{
			List:        conversations,
			HasMore:     hasMore,
			CurrentPage: page,
			PageSize:    pageSize,
		}

		// 3. 返回成功响应
		base.RespJsonSuccess(c, "获取成功", pageResult)
	}
}

// ---------------------- 新增API：分页获取消息 ----------------------
func GetLatestMessagesApiHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		// 1. 获取URL参数
		deviceId := c.GetString(base.ContextDeviceIdKey)
		threadId := r.URL.Query().Get("thread_id")
		pageSizeStr := r.URL.Query().Get("page_size")

		// 2. 参数校验与转换
		if threadId == "" {
			base.RespJsonBadRequest(c, "会话ID不能为空")
			return
		}
		var err error
		pageSize := 20
		if pageSizeStr != "" {
			pageSize, err = strconv.Atoi(pageSizeStr)
			if err != nil || pageSize < 1 || pageSize > 50 { // 限制最大每页50条
				pageSize = 20
			}
		}

		// 3. 查询总消息数（用于计算总页数）
		totalCount, err := getMessageTotalCount(deviceId, threadId)
		if err != nil {
			base.RespJsonInternalServerError(c, "查询消息总数失败："+err.Error())
			return
		}

		// 5. 查询当前页消息
		messages, err := getLatestMessages(deviceId, threadId, pageSize)
		if err != nil {
			base.RespJsonInternalServerError(c, "获取消息列表失败："+err.Error())
			return
		}

		// 判断是否有更多老消息
		hasMore := totalCount > int64(pageSize)

		// 7. 组装分页结果
		pageResult := MessageResult{
			Messages: messages,
			HasMore:  hasMore,
		}

		// 8. 返回成功响应
		base.RespJsonSuccess(c, "获取成功", pageResult)
	}
}

func GetNewMessageApiHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		// 1. 获取URL参数
		deviceId := c.GetString(base.ContextDeviceIdKey)
		threadId := r.URL.Query().Get("thread_id")
		lastDate := r.URL.Query().Get("last_date")
		if threadId == "" || lastDate == "" {
			base.RespJsonBadRequest(c, "参数错误")
			return
		}
		message, err := getNewMessage(deviceId, threadId, lastDate)
		if err != nil {
			base.RespJsonInternalServerError(c, "获取消息失败："+err.Error())
			return
		}
		base.RespJsonSuccess(c, "获取成功", message)
	}
}

// 获取老消息
func GetOldMessagesApiHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		// 获取参数
		deviceId := c.GetString(base.ContextDeviceIdKey)
		threadId := r.URL.Query().Get("thread_id")
		offsetStr := r.URL.Query().Get("offset")
		pageSizeStr := r.URL.Query().Get("page_size")

		// 参数校验
		if threadId == "" {
			base.RespJsonBadRequest(c, "参数错误")
			return
		}

		// 默认参数
		offset := 0
		if offsetStr != "" {
			offset, _ = strconv.Atoi(offsetStr)
			if offset < 0 {
				offset = 0
			}
		}

		pageSize := 20
		if pageSizeStr != "" {
			pageSize, _ = strconv.Atoi(pageSizeStr)
			if pageSize < 1 || pageSize > 50 {
				pageSize = 20
			}
		}

		// 查询总条数
		totalCount, err := getMessageTotalCount(deviceId, threadId)
		if err != nil {
			base.RespJsonInternalServerError(c, "查询消息总数失败："+err.Error())
			return
		}

		// 查询老消息（按date ASC排序，偏移offset，取pageSize条）
		messages, err := getOldMessages(deviceId, threadId, offset, pageSize)
		if err != nil {
			base.RespJsonInternalServerError(c, "获取老消息失败："+err.Error())
			return
		}

		// 判断是否还有更多老消息
		hasMore := (offset + pageSize) < int(totalCount)
		// 组装结果
		result := MessageResult{
			Messages: messages,
			HasMore:  hasMore,
		}

		// 返回响应
		base.RespJsonSuccess(c, "获取成功", result)
	}
}

func SendMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req MessageSendReq
		if err := req.ShouldBindJSON(c, &req); err != nil {
			base.RespJsonBadRequest(c, "请求参数错误")
			return
		}

		if !config.Feature.EnableSendSms {
			base.RespJsonBadRequest(c, "功能未启用")
			return
		}

		dev := c.MustGet(base.ContextDeviceKey).(device.ConnectDevice)

		isms := dev.GetIsms()
		if isms != nil {
			id, err := isms.SendMessage(req.SubId, req.Address, req.Body)
			if err != nil {
				base.RespJsonInternalServerError(c, err.Error())
				return
			}
			base.RespJsonSuccess(c, "发送成功", id)
			return
		}
		base.RespJsonInternalServerError(c, "发送失败")

	}
}
