package sms

import (
	"adb-backup/internal/database"
	"adb-backup/internal/log"
	"time"
)

type Conversation struct {
	ThreadId    string `json:"thread_id"`
	Address     string `json:"address"`
	LastMessage string `json:"last_message"`
	LastTime    string `json:"last_time"`
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
			s1.thread_id,
			s1.address,
			s1.body as last_message,
			s1.date as last_date
		FROM 
			sms s1
		INNER JOIN (
			SELECT 
				thread_id,
				MAX(date) as max_date
			FROM 
				sms
			WHERE 
				device_id = $1
				AND address LIKE $2
			GROUP BY 
				thread_id
		) s2 ON s1.thread_id = s2.thread_id AND s1.date = s2.max_date
		WHERE 
			s1.device_id = $1
		ORDER BY 
			s1.date DESC
		LIMIT $3 OFFSET $4;
	`
	var conversations []Conversation

	rows, err := db.Raw(query, deviceId, "%"+address+"%", pageSize, offset).Rows()
	if err != nil {
		log.WarningF("getConversations error %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Conversation
		var lastDate time.Time
		err := rows.Scan(&c.ThreadId, &c.Address, &c.LastMessage, &lastDate)
		if err != nil {
			return nil, err
		}
		// 格式化时间（ADB获取的date是毫秒级时间戳）
		c.LastTime = lastDate.Format("2006-01-02 15:04")
		conversations = append(conversations, c)
	}
	return conversations, nil
}
