package device

import (
	"adb-backup/internal/config"
	"adb-backup/internal/database"
	"adb-backup/internal/log"
	"adb-backup/internal/notify"
	"adb-backup/internal/shell"
	sy "adb-backup/internal/sync"
	"slices"
	"sync"

	adb "github.com/zach-klippenstein/goadb"

	"time"
)

var (
	client *adb.Adb
	// 设备同步列表
	syncing = make(map[string]*sy.SmsSync)

	connectDevices = make(map[string]*adb.Device)

	devicesMutex = sync.RWMutex{}
)

func GetSyncing() map[string]*sy.SmsSync {
	return syncing
}

func GetConnectDevices() map[string]*adb.Device {
	return connectDevices
}

func GetDevice(serials string) *adb.Device {
	return connectDevices[serials]
}

func GetSmsSync(serials string) *sy.SmsSync {
	return syncing[serials]
}

func StartWatch() {

	initClient()
	// 启动轮询检查设备
	ticker := time.NewTicker(config.App.WaitDeviceInterval)
	defer ticker.Stop()

	log.InfoF("服务已经启动\n开始检测usb设备，请使用usb连接手机并打开<开发者模式>")
	for {
		checkAndSyncDevices()
		<-ticker.C
	}
}

func checkAndSyncDevices() {
	devices, err := client.ListDevices()
	if err != nil {
		log.Fatal(err)
	}
	// 获取当前连接的设备序列号
	var currentSerials []string
	for _, device := range devices {
		currentSerials = append(currentSerials, device.Serial)
	}
	for di := range connectDevices {
		if !slices.Contains(currentSerials, di) {
			delete(connectDevices, di)
		}
	}

	// 处理每个设备
	for _, deviceInfo := range devices {
		handleDevice(deviceInfo)
	}

}

func handleDevice(deviceInfo *adb.DeviceInfo) {
	serial := deviceInfo.Serial

	// 检查设备是否已在同步
	devicesMutex.RLock()
	if _, ok := syncing[serial]; ok {
		devicesMutex.RUnlock()
		return // 已在同步，跳过
	}
	devicesMutex.RUnlock()

	adbDevice := client.Device(adb.DeviceWithSerial(serial))

	connectDevices[serial] = adbDevice

	state, err := adbDevice.State()
	if err != nil {
		log.ErrorF("获取设备 %s 状态错误： %v", deviceInfo.Serial, err)
		return
	}
	if state == adb.StateUnauthorized {
		log.WarningF("设备 %s 未授权，请在打开<USB调试>，并同意USB调试的授权弹框", deviceInfo.Serial)
		return
	}

	if state != adb.StateOnline {
		log.WarningF("设备 %s 状态为 %s，跳过", deviceInfo.Serial, state)
		return
	}
	androidId, aErr := shell.SettingsGetAndroidId(adbDevice)
	if aErr != nil || androidId == "" {
		log.WarningF("设备 %s 无法获取androidId", deviceInfo.Serial)
		return
	}

	// 查找或创建设备记录
	device, err := database.FindDeviceById(androidId)
	if err != nil {
		serialDevice, err := database.FindDeviceById(deviceInfo.Serial)
		if err == nil {
			if serialDevice.Id != "" && serialDevice.Product == deviceInfo.Product && serialDevice.Model == deviceInfo.Model && serialDevice.Info == deviceInfo.DeviceInfo {
				serialDevice.Id = androidId
				database.UpdateDeviceId(serial, androidId)
				database.UpdateSmsDeviceId(serial, androidId)
				device = serialDevice
			}
		}
	}

	if device.Id == "" {
		// 创建新设备
		manufacturer, _ := shell.GetPropProductManufacturer(adbDevice)
		marketingName, _ := shell.GetMarketingName(adbDevice)
		device = database.Device{
			Id:            androidId,
			Serial:        serial,
			Product:       deviceInfo.Product,
			Model:         deviceInfo.Model,
			Info:          deviceInfo.DeviceInfo,
			Usb:           deviceInfo.Usb,
			Manufacturer:  manufacturer,
			MarketingName: marketingName,
		}
		err = database.CreateDevice(&device)
		if err != nil {
			log.FatalF("创建设备 %s 记录错误： %v", serial, err)
			return
		}
	} else {
		var change = false
		if device.Manufacturer == "" {
			manufacturer, _ := shell.GetPropProductManufacturer(adbDevice)
			device.Manufacturer = manufacturer
			change = true
		}
		if device.MarketingName == "" {
			change = true
			marketingName, _ := shell.GetMarketingName(adbDevice)
			device.MarketingName = marketingName
		}
		if change {
			database.UpdateDevice(&device)
		}
	}

	// 标记设备为正在同步并启动同步任务
	devicesMutex.Lock()
	smsSync := sy.SmsSync{
		DbDevice:  device,
		NewNotify: notify.Notify,
		Device:    adbDevice,
	}
	syncing[serial] = &smsSync
	devicesMutex.Unlock()

	// 启动异步同步任务
	go func(dev database.Device) {
		defer func() {
			// 同步完成后清除设备标记
			devicesMutex.Lock()
			delete(syncing, dev.Serial)
			devicesMutex.Unlock()
		}()
		// 执行同步，如果失败则提前返回
		if err := smsSync.StartSync(); err != nil {
			log.ErrorF("设备 %s 同步失败: %v", device.Serial, err)
			return
		}
	}(device)
	log.SuccessF("检测到新的设备 %s，已经开始同步", deviceInfo)
}

func initClient() {
	var err error
	client, err = adb.NewWithConfig(adb.ServerConfig{
		Host: config.Adb.AdbHost,
		Port: config.Adb.AdbPort,
	})
	if err != nil {
		log.FatalF("初始化错误： ", err)
	}
}
