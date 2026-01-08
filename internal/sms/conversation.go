package sms

import (
	"adb-backup/internal/database"
	"time"
)

type Conversation struct {
	ThreadId    string `json:"thread_id"`
	Address     string `json:"address"`
	LastMessage string `json:"last_message"`
	LastTime    string `json:"last_time"`
	SubId       int    `json:"sub_id"`
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

// 获取会话列表的方法（略作修改，返回[]Conversation）
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
