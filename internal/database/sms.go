package database

import (
	"time"
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

func FindLastSms(deviceId string) (Sms, error) {
	var sms Sms
	err := db.Where("device_id = ?", deviceId).Order("date DESC").First(&sms).Error
	return sms, err
}

func CreateInBatches(sms []Sms) error {
	return db.CreateInBatches(sms, 100).Error
}
