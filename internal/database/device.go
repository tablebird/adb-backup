package database

import (
	"adb-backup/internal/utils"
	"time"
)

type Device struct {
	Id            string `gorm:"primaryKey;type:varchar(36);comment:设备ID"`
	Serial        string `gorm:"type:varchar(36);index:serial_idx;comment:设备序列号"`
	Product       string `gorm:"comment:设备产品名称"`
	Model         string `gorm:"comment:设备型号"`
	Info          string `gorm:"comment:设备信息"`
	Usb           string `gorm:"comment:设备USB信息"`
	Manufacturer  string `gorm:"comment:设备制造商"`
	MarketingName string `gorm:"comment:设备营销名称"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (d *Device) BuildName() string {
	var name = utils.AppendPrefix(d.MarketingName, d.Manufacturer, " ")
	if name == "" {
		name = utils.AppendPrefix(d.Model, d.Product, " ")
	}
	return name
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

func FindDeviceById(id string) (Device, error) {
	var device Device
	err := db.Where("id = ?", id).First(&device).Error
	return device, err
}

func CreateDevice(device *Device) error {
	return db.Create(&device).Error
}

func UpdateDevice(device *Device) error {
	return db.Save(&device).Error
}

func UpdateDeviceId(id string, newId string) error {
	return db.Model(&Device{}).Where("id = ?", id).Update("id", newId).Error
}
