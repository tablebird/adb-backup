package database

import (
	"time"
)

type Device struct {
	Id        string `gorm:"primaryKey;type:varchar(36);comment:设备ID"`
	Serial    string `gorm:"type:varchar(36);index:serial_idx;comment:设备序列号"`
	Product   string `gorm:"comment:设备产品名称"`
	Model     string `gorm:"comment:设备型号"`
	Info      string `gorm:"comment:设备信息"`
	Usb       string `gorm:"comment:设备USB信息"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FindAllDevices() ([]Device, error) {
	var devices []Device
	err := db.Find(&devices).Error
	return devices, err
}

func FindDevice(serials []string) ([]Device, error) {
	var devices []Device
	err := db.Where("serial IN ?", serials).Find(&devices).Error
	return devices, err
}

func CreateDevice(device *Device) error {
	return db.Create(&device).Error
}
