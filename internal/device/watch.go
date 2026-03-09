package device

import (
	"adb-backup/internal/config"
	"adb-backup/internal/database"
	"adb-backup/internal/log"
	"sync"

	adb "github.com/zach-klippenstein/goadb"

	"time"
)

var (
	client *adb.Adb

	deviceSerialMap = make(map[string]ConnectDevice)

	idSerialMap = make(map[string]string)

	deviceMutex = sync.RWMutex{}
)

func FindAllDevices() []Device {
	var devices []Device

	for _, value := range deviceSerialMap {
		devices = append(devices, value)
	}
	keys := make([]string, 0, len(idSerialMap))
	for k := range idSerialMap {
		keys = append(keys, k)
	}

	dbs, err := database.FindDeviceByNotInId(keys)
	if err == nil {
		for _, item := range dbs {
			device := newDbDevice(item)
			devices = append(devices, device)
		}
	} else {
		log.WarningF("查询设备失败 %s", err.Error())
	}

	return devices
}

func FindDeviceById(id string) (Device, error) {
	if s, ok := idSerialMap[id]; ok {
		return deviceSerialMap[s], nil
	}
	device, err := database.FindDeviceById(id)
	if err != nil {
		return nil, err
	}
	return newDbDevice(device), nil
}

func StartWatch() {

	client = initClient()

	watcher := client.NewDeviceWatcher()
	// 启动轮询检查设备
	ticker := time.NewTicker(config.App.WaitDeviceInterval)
	defer ticker.Stop()

	log.InfoF("服务已经启动\n开始检测usb设备，请使用usb连接手机并打开<开发者模式>")
	for e := range watcher.C() {
		log.DebugF("设备[%s]  %s > %s", e.Serial, e.OldState, e.NewState)
		if e.CameOnline() {
			device := client.Device(adb.DeviceWithSerial(e.Serial))
			deviceInfo, err := device.DeviceInfo()
			if err == nil {
				handleDevice(deviceInfo)
			} else {
				log.ErrorF("设备[%s] 获取设备信息失败: %v", e.Serial, err)
			}
		} else if e.WentOffline() {
			log.InfoF("设备[%s] 已断开连接", e.Serial)
		}
	}
	watcher.Err()
	log.FatalF("设备监听已关闭[%s]", watcher.Err().Error())
}

func scanAllDevices() {
	devices, err := client.ListDevices()
	if err != nil {
		log.Fatal(err)
	}

	// 处理每个设备
	for _, deviceInfo := range devices {
		handleDevice(deviceInfo)
	}
}

func handleDevice(deviceInfo *adb.DeviceInfo) {
	serial := deviceInfo.Serial

	// 检查设备是否已经存在
	deviceMutex.Lock()
	if phone, ok := deviceSerialMap[serial]; ok {
		initDevice(phone, deviceInfo)
		deviceMutex.Unlock()
		return // 跳过
	}

	adbDevice := client.Device(adb.DeviceWithSerial(serial))
	phone := newConnectDevice(deviceInfo, adbDevice)

	deviceSerialMap[serial] = phone
	initDevice(phone, deviceInfo)
	deviceMutex.Unlock()
}

func initDevice(phone ConnectDevice, deviceInfo *adb.DeviceInfo) {
	oldId := phone.Id()
	delete(idSerialMap, oldId)
	err := phone.initDeviceDB()
	id := phone.Id()
	idSerialMap[id] = deviceInfo.Serial
	if err != nil {
		log.WarningF("保存设备信息失败: %v", err)
		return
	}
	sync := phone.GetSync()
	if sync == nil || sync.StartSync() != nil {
		log.WarningF("启动同步任务失败: %v", err)
		return
	}
	log.SuccessF("检测到新的设备 %s，已经开始同步", phone.Id())
}

func initClient() *adb.Adb {
	client, err := adb.NewWithConfig(adb.ServerConfig{
		Host: config.Adb.AdbHost,
		Port: config.Adb.AdbPort,
	})
	if err != nil {
		log.FatalF("初始化错误： ", err)
	}
	return client
}
