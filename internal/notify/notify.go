package notify

import (
	"adb-backup/internal/database"
	"adb-backup/internal/log"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var Notify Interface

type Interface interface {
	NotifySms(database.Sms) (bool, error)
}

type Webhook struct {
	Url string
}

func (w Webhook) NotifySms(s database.Sms) (bool, error) {
	if s.SmsType != 1 {
		return false, errors.New("not received sms")
	}
	str := fmt.Sprintf(`{"uid": "%s", "address": "%s", "body": "%s"}`, s.Uid, s.Address, s.Body)
	log.DebugF("notifySms : %s", str)
	data := []byte(str)
	resp, err := http.Post(w.Url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.ErrorF("WebHook: %s", err)
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.ErrorF("WebHook: %s", resp.Status)
		return false, errors.New(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.ErrorF("WebHook: %s", err)
		return false, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.ErrorF("WebHook: %s", err)
		return false, err
	}
	uid := result["uid"]

	return uid == s.Uid, nil
}
