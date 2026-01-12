package shell

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	adb "github.com/zach-klippenstein/goadb"
)

func TestServiceCallError(t *testing.T) {
	client, _ := adb.NewWithConfig(adb.ServerConfig{})
	device := client.Device(adb.DeviceWithSerial(DEVICE_TEST_SERIAL))
	res, err := ServiceCall(device, _ISMS, "21")
	assert.NoError(t, err)
	assert.Equal(t, "unknown", res)
	res, err = ServiceCall(device, _ISMS, "1")
	t.Log(res)
	t.Log(err)
}

func TestServiceCallIsmsError(t *testing.T) {
	client, _ := adb.NewWithConfig(adb.ServerConfig{})
	device := client.Device(adb.DeviceWithSerial(DEVICE_TEST_SERIAL))
	res, err := ServiceCall(device, _ISMS, "5",
		_TYPE_INTEGER, strconv.Itoa(0),
		_TYPE_STRING, "com.android.mms.service",
		_TYPE_STRING, _NULL,
		_TYPE_STRING, "10010",
		_TYPE_STRING, _NULL,
		_TYPE_STRING, "message",
		_TYPE_STRING, _NULL,
		_TYPE_STRING, _NULL,
		_TYPE_BOOLEAN, "1",
		_TYPE_BOOLEAN, "0")
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestServiceCallIsmsSendMessage(t *testing.T) {
	client, _ := adb.NewWithConfig(adb.ServerConfig{})
	device := client.Device(adb.DeviceWithSerial(DEVICE_TEST_SERIAL))
	res, err := ServiceCallIsmsSendMessage(device, 0, "10010", "message")
	assert.NoError(t, err)
	assert.True(t, res)
}
