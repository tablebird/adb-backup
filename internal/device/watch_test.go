package device

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitClient(t *testing.T) {
	initClient()
	devices, err := client.ListDevices()
	assert.NoError(t, err)
	for _, item := range devices {
		t.Logf("device %s Serial : %s Product: %s Model: %s", item.DeviceInfo, item.Serial, item.Product, item.Model)
	}
}
