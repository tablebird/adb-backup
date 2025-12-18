package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmsDecode(t *testing.T) {
	data := map[string]string{
		"_id":            "1",
		"thread_id":      "1",
		"date_sent":      "1765007616000",
		"date":           "1765007618426",
		"type":           "1",
		"read":           "1",
		"status":         "-1",
		"subject":        "NULL",
		"service_center": "+12121212",
		"body":           "Hey, are we still meeting at the café at 3 PM today?",
		"address":        "1234567890",
	}
	var sms Sms

	err := sms.Decode(data)
	assert.NoError(t, err)
	assert.Equal(t, 1, sms.Id)
	assert.Equal(t, "1", sms.ThreadId)
	assert.Equal(t, 1, sms.SmsType)
	assert.True(t, sms.Read)
	assert.Equal(t, -1, sms.Status)
	assert.Equal(t, int64(1765007618426), sms.Date.UnixMilli())
	assert.Equal(t, "Hey, are we still meeting at the café at 3 PM today?", sms.Body)
}
