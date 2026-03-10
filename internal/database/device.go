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
	StatusNotify  bool   `gorm:"comment:设备状态通知"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (d *Device) TableName() string {
	return "devices"
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

func FindDeviceByNotInId(ids []string) ([]Device, error) {
	if len(ids) == 0 {
		return FindAllDevices()
	}

	var devices []Device
	err := db.Where("id NOT IN ?", ids).Find(&devices).Error
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

func UpdateDeviceStatusNotify(id string, statusNotify bool) error {
	return db.Model(&Device{}).Where("id = ?", id).Update("status_notify", statusNotify).Error
}
