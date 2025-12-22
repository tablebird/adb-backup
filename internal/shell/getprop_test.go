package shell

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitCommaArray(t *testing.T) {
	value := "value 1,value 2 "
	expected := []string{"value 1", "value 2"}
	actual := strings.Split(strings.TrimSpace(value), ",")
	assert.Equal(t, expected[0], actual[0])
	assert.Equal(t, expected[1], actual[1])

}

func TestSplitCommaArray2(t *testing.T) {
	actual := strings.Split("value 1", ",")
	assert.Equal(t, 1, len(actual))
}
