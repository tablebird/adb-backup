package sync

import (
	"adb-backup/internal/config"
	"adb-backup/internal/database"
	"adb-backup/internal/log"
	"adb-backup/internal/notify"
	"adb-backup/internal/utils"

	"adb-backup/internal/shell"

	"fmt"
	"time"

	uuid "github.com/google/uuid"
	adb "github.com/zach-klippenstein/goadb"
)

type SmsSync struct {
	DbDevice database.Device

	Device *adb.Device

	NewNotify notify.Interface

	startSyncDate time.Time

	smsLastDate time.Time

	contentQuery shell.ContentQuery
}

func (s *SmsSync) SyncSms() error {
	serial := s.DbDevice.Serial
	deviceId := s.DbDevice.Id
	contentQuery := shell.ContentQuery{
		Uri:  shell.CONTENT_QUERY_URI_SMS,
		Sort: "date",
	}
	s.contentQuery = contentQuery
	s.startSyncDate = time.Now()
	lastSms, _ := database.FindLastSms(deviceId)
	s.smsLastDate = lastSms.Date
	log.DebugF("[%s]历史最新消息时间为：%s", serial, s.smsLastDate)

	for {
		if !s.smsLastDate.IsZero() {
			contentQuery.Where = fmt.Sprintf("date>%d", s.smsLastDate.UnixMilli())
		}
		result, err := contentQuery.QueryRow(s.Device)
		if err != nil {
			return fmt.Errorf("读取短信错误： %w", err)
		}
		var messages []database.Sms
		for _, item := range result {
			var sms database.Sms
			utils.MapDecode(shell.ContentQueryParseItem(item), &sms)
			sms.Uid = uuid.New().String()
			sms.DeviceId = deviceId
			sms.OrgStr = utils.CleanString(item)
			if sms.Date.Sub(s.smsLastDate) > 0 {
				s.smsLastDate = sms.Date
			}
			if s.NewNotify != nil && sms.Date.Sub(s.startSyncDate) > 0 {
				s.NewNotify.NotifySms(sms)
			}
			messages = append(messages, sms)
		}

		length := len(messages)
		if length <= 0 {
			wait := config.Conf.ReadInterval
			log.DebugF("[%s]没有找到新短信 暂停%s 最后一条消息的时间为: %s", serial, wait, s.smsLastDate)
			time.Sleep(wait)
		} else {
			err := database.CreateInBatches(messages)
			if err != nil {
				log.WarningF("[%s]保存短信错误： %s", serial, err)
			}
			log.DebugF("[%s]读取到短信数量： %d 最后一条消息的时间为: %s", serial, length, s.smsLastDate)
		}
	}
}
