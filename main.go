package main

import (
	"sync"
	"time"

	adb "github.com/zach-klippenstein/goadb"
	"gorm.io/gorm"
)

var (
	config Config
	client *adb.Adb
	db     *gorm.DB
	notify Interface

	// 同步的设备列表
	syncingDevices = make(map[string]SmsSync)

	devicesMutex = sync.RWMutex{}
)

func main() {
	logInfoF("服务启动中....")

	config.initConfig()
	db = initDB()
	initClient()
	url := config.NotifyWebhookUrl
	if len(url) != 0 {
		notify = Webhook{
			Url: url,
		}
	}

	// 启动轮询检查设备
	ticker := time.NewTicker(config.WaitDeviceInterval)
	defer ticker.Stop()

	logInfoF("服务已经启动\n开始检测usb设备，请使用usb连接手机并打开<开发者模式>")
	for {
		checkAndSyncDevices()
		<-ticker.C
	}
}

func checkAndSyncDevices() {
	devices, err := client.ListDevices()
	if err != nil {
		logFatal(err)
	}
	// 获取当前连接的设备序列号
	var currentSerials []string
	for _, device := range devices {
		currentSerials = append(currentSerials, device.Serial)
	}

	// 查询数据库中的设备
	var dbDevices []Device
	db.Where("serial IN ?", currentSerials).Find(&dbDevices)

	// 处理每个设备
	for _, deviceInfo := range devices {
		handleDevice(deviceInfo, dbDevices)
	}
}

func handleDevice(deviceInfo *adb.DeviceInfo, dbDevices []Device) {
	serial := deviceInfo.Serial

	// 检查设备是否已在同步
	devicesMutex.RLock()
	if _, ok := syncingDevices[serial]; ok {
		devicesMutex.RUnlock()
		return // 已在同步，跳过
	}
	devicesMutex.RUnlock()

	adbDevice := client.Device(adb.DeviceWithSerial(serial))

	state, err := adbDevice.State()
	if err != nil {
		logErrorF("获取设备 %s 状态错误： %v", deviceInfo.Serial, err)
		return
	}
	if state == adb.StateUnauthorized {
		logWarningF("设备 %s 未授权，请在打开<USB调试>，并同意USB调试的授权弹框", deviceInfo.Serial)
		return
	}

	if state != adb.StateOnline {
		logWarningF("设备 %s 状态为 %s，跳过", deviceInfo.Serial, state)
		return
	}

	logSuccessF("检测到新的设备 %s，已经开始同步", deviceInfo)

	// 查找或创建设备记录
	var device Device
	found := false
	for _, dbDevice := range dbDevices {
		if dbDevice.Serial == serial {
			device = dbDevice
			found = true
			break
		}
	}
	if !found {
		// 创建新设备
		device = Device{
			Id:      serial,
			Serial:  serial,
			Product: deviceInfo.Product,
			Model:   deviceInfo.Model,
			Info:    deviceInfo.DeviceInfo,
			Usb:     deviceInfo.Usb,
		}
		db.Create(&device)
	}

	// 标记设备为正在同步并启动同步任务
	devicesMutex.Lock()
	smsSync := SmsSync{
		DbDevice:  device,
		NewNotify: notify,
		Device:    adbDevice,
	}
	syncingDevices[serial] = smsSync
	devicesMutex.Unlock()

	// 启动异步同步任务
	go func(dev Device) {
		defer func() {
			// 同步完成后清除设备标记
			devicesMutex.Lock()
			delete(syncingDevices, dev.Serial)
			devicesMutex.Unlock()
		}()
		// 执行同步，如果失败则提前返回
		if err := smsSync.SyncSms(); err != nil {
			logErrorF("设备 %s 同步失败: %v", device.Serial, err)
			return
		}
	}(device)
}

func initClient() {
	var err error
	client, err = adb.NewWithConfig(adb.ServerConfig{
		Port: config.AdbPort,
	})
	if err != nil {
		logFatalF("初始化错误： ", err)
	}
	err = client.StartServer()

	if err != nil {
		logFatalF("启动adb服务错误： ", err)
	}
}
