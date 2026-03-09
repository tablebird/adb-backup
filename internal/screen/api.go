package screen

import (
	"adb-backup/internal/device"
	"adb-backup/internal/web/base"
	"bufio"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ScreenCap() gin.HandlerFunc {
	return func(c *gin.Context) {
		dev := c.MustGet(base.ContextDeviceKey).(device.ConnectDevice)

		screen := dev.GetScreen()
		if screen == nil {
			base.RespJson(c, http.StatusOK, "设备不支持屏幕截图", nil)
			return
		}

		stdout, err := screen.Capture()
		defer func() {
			if stdout != nil {
				stdout.Close()
			}
		}()
		if err != nil {
			base.RespJsonInternalServerError(c, "获取屏幕失败")
			return
		}
		w := c.Writer
		w.Header().Set("Content-Type", "image/png")
		reader := bufio.NewReader(stdout)
		buf := make([]byte, 4096)
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				if _, writeErr := w.Write(buf[:n]); writeErr != nil {
					return
				}
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
			}
			if err != nil {
				if err == io.EOF {
					return
				}
				return
			}
		}

	}
}

func ScreenView() gin.HandlerFunc {
	return func(c *gin.Context) {
		dev := c.MustGet(base.ContextDeviceKey).(device.ConnectDevice)

		screen := dev.GetScreen()
		if screen == nil {
			base.RespJson(c, http.StatusOK, "设备不支持屏幕直播", nil)
			return
		}

		stdout, err := screen.LiveH264()
		defer func() {
			if stdout != nil {
				stdout.Close()
			}
		}()

		if err != nil {
			base.RespJsonInternalServerError(c, "获取屏幕失败")
			return
		}
		w := c.Writer
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("X-Accel-Buffering", "no")

		reader := bufio.NewReader(stdout)
		buf := make([]byte, 4096)
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				if _, writeErr := w.Write(buf[:n]); writeErr != nil {
					return
				}
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
			}
			if err != nil {
				if err == io.EOF {
					base.RespJsonInternalServerError(c, string(buf[:n]))
					return
				}
				return
			}
		}
	}
}
