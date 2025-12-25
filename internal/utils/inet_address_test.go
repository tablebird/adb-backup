package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocalHostIP(t *testing.T) {
	ip, err := GetLocalHostIP()
	assert.NoError(t, err)
	assert.NotEmpty(t, ip)
	t.Logf("Local Host IP: %s", ip)
}
