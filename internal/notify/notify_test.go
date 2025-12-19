package notify

import (
	"adb-backup/internal/database"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterface(t *testing.T) {
	var i Interface
	assert.Nil(t, i)
}

func TestNotifySms(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		var request map[string]interface{}
		err = json.Unmarshal(body, &request)
		assert.NoError(t, err)
		assert.Equal(t, "1231231231", request["uid"])

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"uid": "1231231231"}`))
	}))
	defer server.Close()

	webhook := Webhook{
		Url: server.URL,
	}
	sms := database.Sms{
		Uid:     "1231231231",
		Address: "911",
		Body:    "test",
	}
	result, err := webhook.NotifySms(sms)
	assert.EqualError(t, err, "not received sms")
	assert.False(t, result)

	sms.SmsType = 1

	result, err = webhook.NotifySms(sms)
	assert.NoError(t, err)
	assert.True(t, result)
}
