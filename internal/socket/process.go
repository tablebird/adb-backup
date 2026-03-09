package socket

import (
	"adb-backup/internal/device"
	"adb-backup/internal/device/touch"
	"adb-backup/internal/log"
	"encoding/json"
	"reflect"

	"github.com/gorilla/websocket"
)

func processText(message []byte, sender DateSender) {
	msg := string(message)
	if msg == "ping" {
		sender.Send(&data{data: []byte("pong"), messageType: websocket.TextMessage})
	} else {
		var req map[string]interface{}
		err := json.Unmarshal(message, &req)
		if err != nil {
			return
		}
		processMap(req, sender)
	}
}

func processMap(req map[string]interface{}, sender DateSender) {
	action, ok := req["a"].(string)
	if !ok {
		log.Warning("action is not string")
		return
	}
	deviceId := req["d"].(string)
	d, err := device.FindDeviceById(deviceId)

	if err != nil {
		log.WarningF("can not find device %s in database", deviceId)
		return
	}
	dev := d.(device.ConnectDevice)
	if dev == nil {
		log.WarningF("device %s is not connect", deviceId)
		return
	}

	switch action {
	case "p":
		input := dev.GetInput()
		if input != nil {
			input.Power()
		}
	case "t":
		input := dev.GetInput()
		if input != nil {
			input.Text(req["t"].(string))
		}
	case "ke":
		log.Debug(reflect.TypeOf(req["kc"]))
		var keyCodes []int
		switch v := req["kc"].(type) {
		case float64: // 单个数值
			keyCodes = []int{int(v)}
		case []interface{}: // 数组
			keyCodes = make([]int, len(v))
			for i, val := range v {
				if fv, ok := val.(float64); ok {
					keyCodes[i] = int(fv)
				} else {
					log.WarningF("keyCode[%d] is not a number", i)
					return
				}
			}
		default:
			log.Warning("keyCode is not a number or array")
			return
		}
		input := dev.GetInput()
		if input != nil {
			input.KeyEvent(keyCodes...)
		}
	case "m":
		touchManager := dev.GetTouch()
		if touchManager != nil {
			ty := req["t"].(float64)
			data := req["p"].(map[string]interface{})
			touchId := data["i"].(float64)
			x := data["x"].(float64)
			y := data["y"].(float64)
			event := touch.NewMouseEvent(int(touchId), touch.MouseType(int(ty)), int(x), int(y))
			if duration, ok := data["d"].(float64); ok {
				event.SetDuration(int(duration))
			}
			touchManager.Mouse(event)
		}
	default:
		// Unknown action
	}
}
