package shell

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScreenCap(t *testing.T) {
	shell := BuildShell()
	r, err := ScreenCap(shell, "-p")
	assert.NoError(t, err)
	defer r.Close()
	buffer := make([]byte, 512)
	n, err := r.Read(buffer)
	assert.NoError(t, err)
	contextType := http.DetectContentType(buffer[:n])
	assert.Equal(t, "image/png", contextType)
}
