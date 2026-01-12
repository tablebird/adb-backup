package shell

import (
	"fmt"
	"strconv"

	adb "github.com/zach-klippenstein/goadb"
)

const (
	_SERVICE = "service"
	_CALL    = "call"
	_ISMS    = "isms"

	_NULL = "null"

	_TYPE_INTEGER = "i32"

	_TYPE_STRING = "s16"

	_TYPE_BOOLEAN = "i32"

	_TYPE_LONG = "i64"
)

func ServiceCall(d *adb.Device, method string, args ...string) (string, error) {
	res, err := d.RunCommand(_SERVICE, append([]string{_CALL, method}, args...)...)
	if err != nil {
		return "", err
	}
	res, err = parseParcel(res)
	if err != nil {
		return "", err
	}
	return res, nil
}

// https://gist.github.com/Ademking/5351ed43a7c48575fe5e6de477d9781f
// https://gist.github.com/lucianogiuseppe/758bb8471fe0119e7b55232343af9ecc
func ServiceCallIsmsSendMessage(d *adb.Device, slot int, phone string, message string) (bool, error) {
	version, err := GetPropBuildVersionRelease(d)
	if err != nil {
		return false, err
	}
	var result string
	var e error

	if version >= 10 {
		result, e = ServiceCall(d, _ISMS, "5",
			_TYPE_INTEGER, strconv.Itoa(slot),
			_TYPE_STRING, "com.android.mms.service",
			_TYPE_STRING, _NULL,
			_TYPE_STRING, phone,
			_TYPE_STRING, _NULL,
			_TYPE_STRING, message,
			_TYPE_STRING, _NULL,
			_TYPE_STRING, _NULL,
			_TYPE_BOOLEAN, "1",
			_TYPE_LONG, "0")
	} else if version == 10 {
		result, e = ServiceCall(d, _ISMS, "7",
			_TYPE_INTEGER, strconv.Itoa(slot),
			_TYPE_STRING, "com.android.mms.service",
			_TYPE_STRING, phone,
			_TYPE_STRING, _NULL,
			_TYPE_STRING, message,
			_TYPE_STRING, _NULL,
			_TYPE_STRING, _NULL,
			_TYPE_BOOLEAN, "1")
	} else if version >= 8 {
		result, e = ServiceCall(d, _ISMS, "7",
			_TYPE_INTEGER, strconv.Itoa(slot),
			_TYPE_STRING, "com.android.mms.service",
			_TYPE_STRING, phone,
			_TYPE_STRING, _NULL,
			_TYPE_STRING, message,
			_TYPE_STRING, _NULL,
			_TYPE_STRING, _NULL)
	} else if version >= 6 {
		result, e = ServiceCall(d, _ISMS, "7",
			_TYPE_INTEGER, strconv.Itoa(slot),
			_TYPE_STRING, "com.android.mms",
			_TYPE_STRING, phone,
			_TYPE_STRING, _NULL,
			_TYPE_STRING, message,
			_TYPE_STRING, _NULL,
			_TYPE_STRING, _NULL)
	} else {
		result, e = ServiceCall(d, _ISMS, "9",
			_TYPE_STRING, "com.android.mms",
			_TYPE_STRING, phone,
			_TYPE_STRING, _NULL,
			_TYPE_STRING, message,
			_TYPE_STRING, _NULL,
			_TYPE_STRING, _NULL)
	}
	if e != nil {
		return false, e
	}
	if result == "" {
		return true, nil
	}
	return false, fmt.Errorf("%s", result)
}
