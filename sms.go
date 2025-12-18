package main

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	uuid "github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	adb "github.com/zach-klippenstein/goadb"
)

type Sms struct {
	Uid           string    `mapstructure:"-" gorm:"primaryKey;type:char(36)"`
	Id            int       `mapstructure:"_id"`
	DeviceId      string    `mapstructure:"-" gorm:"index:device_id_idx,priority:1;not null"`
	ThreadId      string    `mapstructure:"thread_id" gorm:"index:thread_id_idx"`
	Address       string    `mapstructure:"address" gorm:"index:address_idx;index:device_id_idx,priority:2"`
	Date          time.Time `mapstructure:"date"`
	Read          bool      `mapstructure:"read"`
	Status        int       `mapstructure:"status"`
	SmsType       int       `mapstructure:"type"`
	Body          string    `mapstructure:"body" gorm:"type:text"`
	SubId         int       `mapstructure:"sub_id"`
	ServiceCenter string    `mapstructure:"service_center"`
	OrgStr        string    `mapstructure:"-" gorm:"type:text"`
	CreatedAt     time.Time `mapstructure:"-" gorm:"index:device_id_idx,priority:3"`
	UpdatedAt     time.Time `mapstructure:"-"`
}

func (s *Sms) Decode(data map[string]string) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: stringToTimeHookFunc(),
		Result:     s,
	})
	if err != nil {
		return err
	}
	err = decoder.Decode(data)
	return err
}

func stringToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		str := data.(string)
		switch t.Kind() {
		case reflect.Struct:
			if t == reflect.TypeOf(time.Time{}) {
				if value, err := strconv.ParseInt(str, 10, 64); err == nil {
					return time.UnixMilli(value), nil
				}
			}
		case reflect.Bool:
			return str == "1" || str == "true", nil
		case reflect.Int:
			if val, err := strconv.Atoi(str); err == nil {
				return val, nil
			}
		}
		return data, nil
	}
}

type SmsSync struct {
	DbDevice Device

	Device *adb.Device

	NewNotify Interface

	startSyncDate time.Time

	smsLastDate time.Time

	contentQuery ContentQuery
}

func (s *SmsSync) SyncSms() error {
	serial := s.DbDevice.Serial
	deviceId := s.DbDevice.Id
	contentQuery := ContentQuery{
		uri:  uriSms,
		sort: "date",
	}
	s.contentQuery = contentQuery
	s.startSyncDate = time.Now()
	var sms []Sms
	db.Order("date desc").Limit(1).Find(&sms)
	if len(sms) > 0 {
		s.smsLastDate = sms[0].Date
	}
	logDebugF("[%s]历史最新消息时间为：%s", serial, s.smsLastDate)

	for {
		if !s.smsLastDate.IsZero() {
			contentQuery.where = fmt.Sprintf("date>%d", s.smsLastDate.UnixMilli())
		}
		result, err := contentQuery.QueryRow(s.Device)
		if err != nil {
			return fmt.Errorf("读取短信错误： %w", err)
		}
		var messages []Sms
		for _, item := range result {
			var sms Sms
			sms.Decode(parseItem(item))
			sms.Uid = uuid.New().String()
			sms.DeviceId = deviceId
			sms.OrgStr = cleanString(item)
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
			wait := config.ReadInterval
			logDebugF("[%s]没有找到新短信 暂停%s 最后一条消息的时间为: %s", serial, wait, s.smsLastDate)
			time.Sleep(wait)
		} else {
			result := db.CreateInBatches(messages, 100)
			if result.Error != nil {
				logWarningF("[%s]保存短信错误： %s", serial, result.Error)
			}
			logDebugF("[%s]读取到短信数量： %d 最后一条消息的时间为: %s", serial, length, s.smsLastDate)
		}
	}
}
