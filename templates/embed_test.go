package templates

import (
    "testing"
)

func TestGetFiles(t *testing.T) {
    files := getFiles()
    if len(files) == 0 {
        t.Errorf("getFiles() = %v, want %v", files, []string{"index.html"})
    }
    t.Logf("files %v", files)
}