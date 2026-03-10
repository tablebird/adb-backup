package notify

import (
	"adb-backup/internal/config"
	"adb-backup/internal/database"
	"adb-backup/internal/log"
	"errors"
	"fmt"
)

var notify Interface

func GetNotify() Interface {
	return notify
}

type Interface interface {
	NotifySms(database.Sms) (bool, error)

	NotifyDeviceStatus(dev *database.Device, status string, ty string) (bool, error)
}

func InitNotify() {
	conf := config.Notify
	if conf.NotifyWebhookUrl == "" && conf.NotifyStatusWebhookUrl == "" {
		log.Info("未配置通知功能")
		return
	}
	notify = webhook{
		Url:             conf.NotifyWebhookUrl,
		DeviceStatusUrl: conf.NotifyStatusWebhookUrl,
	}
}

type webhook struct {
	Url             string
	DeviceStatusUrl string
}

func (w webhook) NotifySms(s database.Sms) (bool, error) {
	if s.SmsType != 1 {
		return false, errors.New("not received sms")
	}
	if len(w.Url) == 0 {
		return false, errors.New("notify webhook url is empty")
	}
	str := fmt.Sprintf(`{"uid": "%s", "address": "%s", "body": "%s"}`, s.Uid, s.Address, s.Body)
	log.DebugF("notifySms : %s", str)
	result, err := postWebhookJsonStr(w.Url, str)
	if err != nil {
		log.ErrorF("WebHook: %s", err)
		return false, err
	}
	uid := result["uid"]

	return uid == s.Uid, nil
}

func (w webhook) NotifyDeviceStatus(device *database.Device, status string, ty string) (bool, error) {
	if len(w.DeviceStatusUrl) == 0 {
		return false, errors.New("notify status webhook url is empty")
	}
	str := fmt.Sprintf(`{"id": "%s", "serial": "%s", "name": "%s", "status": "%s", "type": "%s"}`, device.Id, device.Serial, device.BuildName(), status, ty)
	log.DebugF("notifyDeviceStatus : %s", str)
	_, err := postWebhookJsonStr(w.DeviceStatusUrl, str)
	return err == nil, err
}
