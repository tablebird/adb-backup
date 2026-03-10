package notify

import (
	"adb-backup/internal/log"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func postWebhookJsonStr(url string, jsonStr string) (map[string]interface{}, error) {
	data := []byte(jsonStr)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.ErrorF("WebHook: %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.ErrorF("WebHook: %s", resp.Status)
		return nil, errors.New(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.ErrorF("WebHook: %s", err)
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	return result, err
}
