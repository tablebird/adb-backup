package sms

import (
	"adb-backup/internal/database"
	"adb-backup/internal/log"

	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Conversation struct {
	ThreadId    string `json:"thread_id"`
	Address     string `json:"address"`
	LastMessage string `json:"last_message"`
	LastTime    string `json:"last_time"`
	SubId       int    `json:"sub_id"`
}

type ConversationPageResult struct {
	List        []Conversation `json:"list"`         // 当前页会话列表
	HasMore     bool           `json:"has_more"`     // 是否还有更多会话
	CurrentPage int            `json:"current_page"` // 当前页码
	PageSize    int            `json:"page_size"`    // 每页条数
}

type Message struct {
	ThreadId   string    `json:"thread_id"`
	Address    string    `json:"address"`
	Date       time.Time `json:"date"`
	DateFormat string    `json:"date_format"`
	SmsType    string    `json:"sms_type"`
	Body       string    `json:"body"`
	SubId      int       `json:"sub_id"`
}

type MessageResult struct {
	Messages []Message `json:"messages"` // 最新消息列表（已反转，最新在最后）
	HasMore  bool      `json:"has_more"` // 是否有更多老消息
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

// 查询设备的会话总数
func getConversationTotalCount(deviceId string, address string) (int64, error) {
	db := database.GetDB()

	query := `SELECT COUNT(DISTINCT thread_id) FROM sms WHERE device_id = $1 AND address LIKE $2;`
	var totalCount int64
	err := db.Raw(query, deviceId, "%"+address+"%").Scan(&totalCount).Error
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

// 原有获取会话列表的方法（略作修改，返回[]Conversation）
func getConversations(deviceId string, offset int, pageSize int, address string) ([]Conversation, error) {
	db := database.GetDB()

	// 注：假设你的短信表名为sms，且有device_id字段关联设备
	query := `
		SELECT 
			thread_id, 
			address, 
			(SELECT body FROM sms WHERE thread_id = t.thread_id AND device_id = $1 ORDER BY date DESC LIMIT 1) as last_message,
			(SELECT date FROM sms WHERE thread_id = t.thread_id AND device_id = $1 ORDER BY date DESC LIMIT 1) as last_date,
			(SELECT sub_id FROM sms WHERE thread_id = t.thread_id AND device_id = $1 ORDER BY date DESC LIMIT 1) as sub_id
		FROM 
			sms t
		WHERE 
			t.device_id = $1
			AND t.address LIKE $2
		GROUP BY 
			t.thread_id, t.address
		ORDER BY 
			last_date DESC
		LIMIT $3 OFFSET $4;
	`
	var conversations []Conversation

	rows, err := db.Raw(query, deviceId, "%"+address+"%", pageSize, offset).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Conversation
		var lastDate time.Time
		err := rows.Scan(&c.ThreadId, &c.Address, &c.LastMessage, &lastDate, &c.SubId)
		if err != nil {
			return nil, err
		}
		// 格式化时间（ADB获取的date是毫秒级时间戳）
		c.LastTime = lastDate.Format("2006-01-02 15:04")
		conversations = append(conversations, c)
	}
	return conversations, nil
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

// 查询指定会话的消息总数
func getMessageTotalCount(deviceId string, threadId string) (int64, error) {
	db := database.GetDB()

	query := `SELECT COUNT(*) FROM sms WHERE device_id = ? AND thread_id = ?;`
	var totalCount int64
	err := db.Raw(query, deviceId, threadId).Scan(&totalCount).Error
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

// 分页查询指定会话的消息
func getLatestMessages(deviceId string, threadId string, limit int) ([]Message, error) {
	db := database.GetDB()

	query := `
		SELECT 
			thread_id, address, date, sms_type, body, sub_id
		FROM 
			sms
		WHERE 
			device_id = ? AND thread_id = ?
		ORDER BY 
			date DESC
		LIMIT ?;
	`

	rows, err := db.Raw(query, deviceId, threadId, limit).Rows()
	if err != nil {
		log.ErrorF("查询消息列表失败：%v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		err := rows.Scan(&m.ThreadId, &m.Address, &m.Date, &m.SmsType, &m.Body, &m.SubId)
		if err != nil {
			log.ErrorF("扫描消息行失败：%v", err)
			return nil, err
		}
		// 格式化时间
		m.DateFormat = m.Date.Format("2006-01-02 15:04:05")
		messages = append(messages, m)
	}
	// 反转消息列表（让最新的消息在数组最后，对应页面底部）
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages, nil
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

// 查询老消息
func getOldMessages(deviceId string, threadId string, offset, pageSize int) ([]Message, error) {
	db := database.GetDB()

	query := `
        SELECT 
            thread_id, address, date, sms_type, body, sub_id
        FROM 
            sms
        WHERE 
            device_id = $1 AND thread_id = $2
        ORDER BY 
            date ASC
        LIMIT $3 OFFSET $4;
    `

	rows, err := db.Raw(query, deviceId, threadId, pageSize, offset).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		err := rows.Scan(&m.ThreadId, &m.Address, &m.Date, &m.SmsType, &m.Body, &m.SubId)
		if err != nil {
			return nil, err
		}
		m.DateFormat = m.Date.Format("2006-01-02 15:04:05")
		messages = append(messages, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}
