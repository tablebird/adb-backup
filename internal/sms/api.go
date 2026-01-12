package sms

import (
	"adb-backup/internal/device"
	"adb-backup/internal/shell"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	adb "github.com/zach-klippenstein/goadb"
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
	DeviceId string `json:"device_id" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Body     string `json:"body" binding:"required"`
	SubId    int    `json:"sub_id" binding:"oneof=0 1"`
}

func GetConversationsApiHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		// 1. 获取设备ID参数
		deviceId := r.URL.Query().Get("device_id")
		pageStr := r.URL.Query().Get("page")
		pageSizeStr := r.URL.Query().Get("page_size")
		address := r.URL.Query().Get("address")

		if deviceId == "" {
			// 返回错误响应
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "设备ID不能为空",
				"data": nil,
			})
			return
		}
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
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "查询会话总数失败：" + err.Error(),
				"data": nil,
			})
			return
		}

		// 计算偏移量
		offset := (page - 1) * pageSize
		// 2. 查询会话列表
		conversations, err := getConversations(deviceId, offset, pageSize, address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "获取会话列表失败：" + err.Error(),
				"data": nil,
			})
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
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取成功",
			"data": pageResult,
		})
	}
}

// ---------------------- 新增API：分页获取消息 ----------------------
func GetLatestMessagesApiHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		// 1. 获取URL参数
		deviceId := r.URL.Query().Get("device_id")
		threadId := r.URL.Query().Get("thread_id")
		pageSizeStr := r.URL.Query().Get("page_size")

		// 2. 参数校验与转换
		if deviceId == "" || threadId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "设备ID和会话ID不能为空",
				"data": nil,
			})
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
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "查询消息总数失败：" + err.Error(),
				"data": nil,
			})
			return
		}

		// 5. 查询当前页消息
		messages, err := getLatestMessages(deviceId, threadId, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "获取消息列表失败：" + err.Error(),
				"data": nil,
			})
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
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取成功",
			"data": pageResult,
		})
	}
}

func GetNewMessageApiHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		// 1. 获取URL参数
		deviceId := r.URL.Query().Get("device_id")
		threadId := r.URL.Query().Get("thread_id")
		lastDate := r.URL.Query().Get("last_date")
		if deviceId == "" || threadId == "" || lastDate == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "参数错误",
				"data": nil,
			})
			return
		}
		message, err := getNewMessage(deviceId, threadId, lastDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "获取消息失败：" + err.Error(),
				"data": nil,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取成功",
			"data": message,
		})
	}
}

// 获取老消息
func GetOldMessagesApiHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		// 获取参数
		deviceId := r.URL.Query().Get("device_id")
		threadId := r.URL.Query().Get("thread_id")
		offsetStr := r.URL.Query().Get("offset")
		pageSizeStr := r.URL.Query().Get("page_size")

		// 参数校验
		if deviceId == "" || threadId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "设备ID和会话ID不能为空",
				"data": nil,
			})
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
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "查询消息总数失败：" + err.Error(),
				"data": nil,
			})
			return
		}

		// 查询老消息（按date ASC排序，偏移offset，取pageSize条）
		messages, err := getOldMessages(deviceId, threadId, offset, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "获取老消息失败：" + err.Error(),
				"data": nil,
			})
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
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取成功",
			"data": result,
		})
	}
}

func SendMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req MessageSendReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "请求参数错误",
				"data": err.Error(),
			})
			return
		}

		dev := device.GetDevice(req.DeviceId)
		if dev == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "设备不存在",
				"data": nil,
			})
			return
		}

		state, er := dev.State()
		if er != nil || state != adb.StateOnline {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "设备状态异常",
				"data": nil,
			})
			return
		}
		networkTypes, networkTypeErr := shell.GetPropGsmNetworkType(dev)
		if networkTypeErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "系统内部错误",
				"data": nil,
			})
			return
		}

		if req.SubId >= len(networkTypes) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "无效的SIM卡",
				"data": nil,
			})
			return
		}

		networkType := networkTypes[req.SubId]

		if len(networkType) == 0 || strings.ToUpper(networkType) == "UNKNOWN" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "Sim卡不可用",
				"data": nil,
			})
			return
		}

		res, sendErr := shell.ServiceCallIsmsSendMessage(dev, req.SubId, req.Address, req.Body)
		if sendErr != nil || !res {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "发送失败",
				"data": nil,
			})
		}

		sync := device.GetSmsSync(req.DeviceId)
		if sync != nil {
			sync.SyncSms()
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "发送成功",
			"data": nil,
		})

	}
}
