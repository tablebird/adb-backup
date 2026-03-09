package shell

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScreenRecord(t *testing.T) {
	device := BuildShell()
	stdout, err := ScreenRecord(device, "--output-format=h264", "--size=1280x720", "--time-limit=1", "-")
	assert.NoError(t, err)

	defer stdout.Close()

	buffer := make([]byte, 4)
	_, err = io.ReadFull(stdout, buffer)
	assert.NoError(t, err)
	// simple verification of h264 stream
	assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x01}, buffer)
}

func TestScreenRecordH264Live(t *testing.T) {
	device := BuildShell()
	stdout, err := ScreenRecordH264Live(device)
	assert.NoError(t, err)

	defer stdout.Close()

	var i = 0
	buffer := make([]byte, 1024*32)
	for {
		i++
		_, err := stdout.Read(buffer)
		if err == io.EOF {
			assert.NoError(t, err)
			break
		}
		if i > 10 {
			break
		}
	}
	assert.Equal(t, 11, i)
}
