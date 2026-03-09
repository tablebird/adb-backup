package isms

import (
	"adb-backup/internal/shell"
	"adb-backup/internal/sync"
	"errors"
	"strings"
)

type IsmsManager interface {
	SendMessage(subId int, address string, body string) (string, error)
}

func NewIsmsManager(s shell.Shell, smsSync sync.SmsSync) IsmsManager {
	return &shellIsms{
		s:       s,
		smsSync: smsSync,
	}
}

type shellIsms struct {
	s       shell.Shell
	smsSync sync.SmsSync
}

func (i *shellIsms) SendMessage(subId int, address string, body string) (id string, resErr error) {
	networkTypes, err := shell.GetPropGsmNetworkType(i.s)

	if err != nil {
		resErr = err
		return
	}
	if subId >= len(networkTypes) {
		resErr = errors.New("无效的SIM卡")
		return
	}

	networkType := networkTypes[subId]
	if len(networkType) == 0 || strings.ToUpper(networkType) == "UNKNOWN" {
		resErr = errors.New("Sim卡不可用")
		return
	}
	res, err := shell.ServiceCallIsmsSendMessage(i.s, subId, address, body)
	if err != nil {
		resErr = err
		return
	}
	if !res {
		resErr = errors.New("发送失败")
		return
	}
	if i.smsSync != nil {
		messages, err := i.smsSync.SyncNow()
		if err != nil {
			resErr = err
			return
		}
		if len(messages) > 0 {
			for _, message := range messages {
				if message.Address == address && message.Body == body {
					id = message.ThreadId
					return
				}
			}
		}

	}
	resErr = errors.New("发送结果未知")
	return
}
