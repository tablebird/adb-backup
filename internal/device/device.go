package device

import (
	"adb-backup/internal/database"
	"adb-backup/internal/log"
	"adb-backup/internal/notify"
	"adb-backup/internal/shell"
	"adb-backup/internal/sync"
	"fmt"

	adb "github.com/zach-klippenstein/goadb"
)

type Device interface {
	Id() string

	Name() string

	State() DeviceState
}

type ConnectDevice interface {
	Device

	initDeviceDB() error

	GetSync() sync.Sync

	GetTelephony() TelephonyManager

	GetWifi() WifiManager

	GetPower() PowerManager

	GetBuild() Build

	GetIsms() Isms
}

func newDbDevice(device database.Device) Device {
	return &dbDevice{deviceDB: &device}
}

func newConnectDevice(deviceInfo *adb.DeviceInfo, adbDevice *adb.Device) ConnectDevice {
	return &shellConnectDevice{
		deviceInfo: deviceInfo,
		adbDevice:  adbDevice}
}

type dbDevice struct {
	deviceDB *database.Device
}

func (p *dbDevice) Id() string {
	return p.deviceDB.Id
}

func (p *dbDevice) Name() string {
	return p.deviceDB.BuildName()
}

func (p *dbDevice) State() DeviceState {
	return StateDisconnected
}

type shellConnectDevice struct {
	dbDevice

	deviceInfo *adb.DeviceInfo

	adbDevice *adb.Device

	sync sync.Sync
}

func (p *shellConnectDevice) Id() string {
	if p.deviceDB == nil {
		return p.deviceInfo.Serial
	}
	return p.dbDevice.Id()
}

func (p *shellConnectDevice) Name() string {
	if p.deviceDB == nil {
		return p.deviceInfo.Model
	}
	return p.dbDevice.Name()
}

func (p *shellConnectDevice) State() DeviceState {
	state, err := p.adbDevice.State()
	if err != nil {
		return StateError
	}
	return deviceStateToStr(state)
}

func (p *shellConnectDevice) GetTelephony() TelephonyManager {
	if p.State() != StateOnline {
		return nil
	}
	return &shellTelephony{adbDevice: p.adbDevice}
}

func (p *shellConnectDevice) GetWifi() WifiManager {
	if p.State() != StateOnline {
		return nil
	}
	return &shellWifi{adbDevice: p.adbDevice}
}

func (p *shellConnectDevice) GetPower() PowerManager {
	if p.State() != StateOnline {
		return nil
	}
	return &shellPower{adbDevice: p.adbDevice}
}

func (p *shellConnectDevice) GetBuild() Build {
	if p.State() != StateOnline {
		return nil
	}
	return &shellBuild{adbDevice: p.adbDevice}
}

func (p *shellConnectDevice) GetIsms() Isms {
	if p.State() != StateOnline {
		return nil
	}
	return &shellIsms{adbDevice: p.adbDevice, smsSync: p.sync.(sync.SmsSync)}
}

func (p *shellConnectDevice) initDeviceDB() error {
	state := p.State()

	if state == StateUnauthorized {
		return fmt.Errorf("设备 %s 未授权，请在打开<USB调试>，并同意USB调试的授权弹框", p.deviceInfo.Serial)
	}

	if state != StateOnline {
		return fmt.Errorf("设备 %s 状态为 %s，跳过", p.deviceInfo.Serial, state)
	}
	androidId, aErr := shell.SettingsGetAndroidId(p.adbDevice)
	if aErr != nil || androidId == "" {
		return fmt.Errorf("设备 %s 无法获取androidId", p.deviceInfo.Serial)
	}

	// 序列号重复导致的设备不匹配
	if p.deviceDB != nil && p.Id() == androidId {
		return nil
	}
	p.deviceDB = nil

	// 查找或创建设备记录
	device, err := database.FindDeviceById(androidId)
	if err != nil {
		serialDevice, err := database.FindDeviceById(p.deviceInfo.Serial)
		if err == nil {
			if serialDevice.Id != "" && serialDevice.Product == p.deviceInfo.Product && serialDevice.Model == p.deviceInfo.Model && serialDevice.Info == p.deviceInfo.DeviceInfo {
				serialDevice.Id = androidId
				database.UpdateDeviceId(p.deviceInfo.Serial, androidId)
				database.UpdateSmsDeviceId(p.deviceInfo.Serial, androidId)
				device = serialDevice
			}
		}
	}

	if device.Id == "" {
		// 创建新设备
		manufacturer, _ := shell.GetPropProductManufacturer(p.adbDevice)
		marketingName, _ := shell.GetMarketingName(p.adbDevice)
		device = database.Device{
			Id:            androidId,
			Serial:        p.deviceInfo.Serial,
			Product:       p.deviceInfo.Product,
			Model:         p.deviceInfo.Model,
			Info:          p.deviceInfo.DeviceInfo,
			Usb:           p.deviceInfo.Usb,
			Manufacturer:  manufacturer,
			MarketingName: marketingName,
		}
		err = database.CreateDevice(&device)
		if err != nil {
			return fmt.Errorf("创建设备 %s 记录错误： %v", p.deviceInfo.Serial, err)
		}
	} else {
		var change = false
		if device.Manufacturer == "" {
			manufacturer, _ := shell.GetPropProductManufacturer(p.adbDevice)
			device.Manufacturer = manufacturer
			change = true
		}
		if device.MarketingName == "" {
			change = true
			marketingName, _ := shell.GetMarketingName(p.adbDevice)
			device.MarketingName = marketingName
		}
		if change {
			database.UpdateDevice(&device)
		}
	}
	p.deviceDB = &device
	return nil
}

func (p *shellConnectDevice) startSync() error {
	if p.deviceDB == nil {
		return fmt.Errorf("设备未保存，无法同步")
	}
	if p.sync != nil {
		return nil
	}
	p.sync = sync.NewSmsSync(p.deviceDB, p.adbDevice, notify.Notify)
	go func() {
		defer func() {
			p.sync = nil
		}()
		if err := p.sync.StartSync(); err != nil {
			log.ErrorF("设备[%s] 同步失败: %v", p.deviceDB.Serial, err)
		}
	}()
	return nil
}

func (p *shellConnectDevice) GetSync() sync.Sync {
	if p.sync != nil {
		return p.sync
	}
	if p.deviceDB == nil {
		return nil
	}
	if p.State() != StateOnline {
		return nil
	}
	p.sync = sync.NewSmsSync(p.deviceDB, p.adbDevice, notify.Notify)
	return p.sync
}
