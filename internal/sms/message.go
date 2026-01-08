package sms

import (
	"adb-backup/internal/database"
	"adb-backup/internal/log"
	"time"
)

type Message struct {
	ThreadId   string    `json:"thread_id"`
	Address    string    `json:"address"`
	Date       time.Time `json:"date"`
	DateFormat string    `json:"date_format"`
	SmsType    string    `json:"sms_type"`
	Body       string    `json:"body"`
	SubId      int       `json:"sub_id"`
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
