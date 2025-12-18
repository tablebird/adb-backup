package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Interface interface {
	NotifySms(Sms) (bool, error)
}

type Webhook struct {
	Url string
}

func (w Webhook) NotifySms(s Sms) (bool, error) {
	if s.SmsType != 1 {
		return false, errors.New("not received sms")
	}
	str := fmt.Sprintf(`{"uid": "%s", "address": "%s", "body": "%s"}`, s.Uid, s.Address, s.Body)
	logDebugF("notifySms : %s", str)
	data := []byte(str)
	resp, err := http.Post(w.Url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		logErrorF("WebHook: %s", err)
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logErrorF("WebHook: %s", resp.Status)
		return false, errors.New(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logErrorF("WebHook: %s", err)
		return false, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		logErrorF("WebHook: %s", err)
		return false, err
	}
	uid := result["uid"]

	return uid == s.Uid, nil
}
