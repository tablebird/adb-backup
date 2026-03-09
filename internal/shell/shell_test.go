package shell

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	adb "github.com/tablebird/goadb"
)

const (
	DEVICE_TEST_SERIAL = "emulator-5554"
)

func BuildShell() AnyShell {
	client, _ := adb.NewWithConfig(adb.ServerConfig{})
	return NewShell(client, DEVICE_TEST_SERIAL)
}

func TestAsyncShellDialDevice(t *testing.T) {
	client, _ := adb.NewWithConfig(adb.ServerConfig{})
	shell := NewShell(client, DEVICE_TEST_SERIAL).(*realShell)
	conn, err := shell.dialDevice()
	assert.Nil(t, err)
	assert.NotNil(t, conn)
	defer conn.Close()
	err = conn.SendMessage([]byte("shell:"))
	assert.Nil(t, err)
	_, err = conn.ReadStatus("shell:")
	assert.Nil(t, err)
	cmd := "echo test\n"
	i, err := conn.Write([]byte(cmd))
	assert.Nil(t, err)
	assert.Equal(t, i, len(cmd))
	time.Sleep(time.Millisecond * 100)
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	assert.Nil(t, err)
	assert.NotEqual(t, n, 0)
	assert.NotEmpty(t, string(buf[:n]))
}

func TestAsyncShellAsyncRunCommand(t *testing.T) {
	client, _ := adb.NewWithConfig(adb.ServerConfig{})
	shell := NewShell(client, DEVICE_TEST_SERIAL).(*realShell)
	shell.AsyncRunCommand("echo 1")
	shell.AsyncRunCommand("echo 2")
	shell.AsyncRunCommand("echo 3")
	for i := 0; i < 9; i++ {
		assert.NotEqual(t, "", <-shell.out)
	}
}
