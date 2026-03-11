package device

import (
	"adb-backup/internal/database"
	"adb-backup/internal/device/build"
	"adb-backup/internal/device/display"
	"adb-backup/internal/device/input"
	"adb-backup/internal/device/isms"
	"adb-backup/internal/device/power"
	"adb-backup/internal/device/screen"
	"adb-backup/internal/device/telephony"
	"adb-backup/internal/device/touch"
	"adb-backup/internal/device/wifi"
	"adb-backup/internal/notify"
	"adb-backup/internal/shell"
	"adb-backup/internal/sync"
	"fmt"

	adb "github.com/tablebird/goadb"
)

type Device interface {
	Id() string

	Name() string

	State() DeviceState

	GetDeviceDB() *database.Device
}

type ConnectDevice interface {
	Device

	initDeviceDB() error

	GetSync() sync.Sync

	GetTelephony() telephony.TelephonyManager

	GetWifi() wifi.WifiManager

	GetPower() power.PowerManager

	GetBuild() build.Build

	GetIsms() isms.IsmsManager

	GetDisplay() display.DisplayManager

	GetTouch() touch.TouchManager

	GetScreen() screen.ScreenManager

	GetInput() input.InputManager
}

func newDbDevice(device database.Device) Device {
	return &dbDevice{deviceDB: &device}
}

func newConnectDevice(deviceInfo *adb.DeviceInfo, adbDevice *adb.Device) ConnectDevice {
	s := shell.NewShell(client, deviceInfo.Serial)
	return &shellConnectDevice{
		deviceInfo: deviceInfo,
		adbDevice:  adbDevice,
		shell:      s}
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

func (p *dbDevice) GetDeviceDB() *database.Device {
	return p.deviceDB
}

type shellConnectDevice struct {
	dbDevice

	deviceInfo *adb.DeviceInfo

	adbDevice *adb.Device

	shell shell.AnyShell

	sync sync.Sync

	touch touch.TouchManager

	screen screen.ScreenManager
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

func (p *shellConnectDevice) GetTelephony() telephony.TelephonyManager {
	if p.State() != StateOnline {
		return nil
	}
	return telephony.NewTelephonyManager(p.shell)
}

func (p *shellConnectDevice) GetWifi() wifi.WifiManager {
	if p.State() != StateOnline {
		return nil
	}
	return wifi.NewWifiManager(p.shell)
}

func (p *shellConnectDevice) GetPower() power.PowerManager {
	if p.State() != StateOnline {
		return nil
	}
	return power.NewPowerManager(p.shell)
}

func (p *shellConnectDevice) GetBuild() build.Build {
	if p.State() != StateOnline {
		return nil
	}
	return build.NewBuild(p.shell)
}

func (p *shellConnectDevice) GetIsms() isms.IsmsManager {
	if p.State() != StateOnline {
		return nil
	}
	return isms.NewIsmsManager(p.shell, p.sync.(sync.SmsSync))
}

func (p *shellConnectDevice) GetDisplay() display.DisplayManager {
	if p.State() != StateOnline {
		return nil
	}
	return display.NewDisplayManager(p.shell)
}

func (p *shellConnectDevice) GetTouch() touch.TouchManager {
	if p.State() != StateOnline {
		return nil
	}
	if p.touch == nil {
		p.touch = touch.NewTouchManager(p.shell)
	}
	return p.touch
}

func (p *shellConnectDevice) GetScreen() screen.ScreenManager {
	if p.State() != StateOnline {
		return nil
	}
	p.screen = screen.NewScreenManager(p.shell)
	return p.screen
}

func (p *shellConnectDevice) GetInput() input.InputManager {
	if p.State() != StateOnline {
		return nil
	}

	return input.NewInputManager(p.shell)
}

func (p *shellConnectDevice) initDeviceDB() error {
	state := p.State()

	if state == StateUnauthorized {
		return fmt.Errorf("设备 %s 未授权，请在打开<USB调试>，并同意USB调试的授权弹框", p.deviceInfo.Serial)
	}

	if state != StateOnline {
		return fmt.Errorf("设备 %s 状态为 %s，跳过", p.deviceInfo.Serial, state)
	}
	androidId, aErr := shell.SettingsGetAndroidId(p.shell)
	if aErr != nil || androidId == "" {
		return fmt.Errorf("设备 %s 无法获取androidId %s", p.deviceInfo.Serial, aErr.Error())
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
		manufacturer, _ := shell.GetPropProductManufacturer(p.shell)
		marketingName, _ := shell.GetMarketingName(p.shell)
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
			manufacturer, _ := shell.GetPropProductManufacturer(p.shell)
			device.Manufacturer = manufacturer
			change = true
		}
		if device.MarketingName == "" {
			change = true
			marketingName, _ := shell.GetMarketingName(p.shell)
			device.MarketingName = marketingName
		}
		if change {
			database.UpdateDevice(&device)
		}
	}
	p.deviceDB = &device
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
	p.sync = sync.NewSmsSync(p.deviceDB, p.shell, notify.GetNotify())
	return p.sync
}
