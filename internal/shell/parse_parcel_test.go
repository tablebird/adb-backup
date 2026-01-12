package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseParcel(t *testing.T) {
	text := `Result: Parcel(    00000000    '....')`
	res, err := parseParcel(text)
	assert.NoError(t, err)
	assert.Equal(t, "", res)
}

func TestParseParcelError(t *testing.T) {
	text := `xxxxx error`
	res, err := parseParcel(text)
	assert.Error(t, err)
	assert.Equal(t, "", res)
}
