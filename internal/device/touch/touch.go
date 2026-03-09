package touch

import (
	"adb-backup/internal/shell"
	"sync/atomic"
)

type TouchManager interface {
	Size() (int, int)

	Mouse(event *MouseEvent)
}

func NewTouchManager(s shell.Shell) TouchManager {
	device := newEventDevice(s)
	return &mouseTouch{device: device,
		mouseRun: 0, mouseEvent: make(chan *MouseEvent, 100)}
}

type mouseTouch struct {
	device mouseDevice

	mouseEvent chan *MouseEvent
	mouseRun   int32
}

func (m *mouseTouch) Size() (int, int) {
	return m.device.size()
}

func (m *mouseTouch) Mouse(event *MouseEvent) {
	go m.startMouse()
	m.mouseEvent <- event
}

func (m *mouseTouch) startMouse() {
	if !atomic.CompareAndSwapInt32(&m.mouseRun, 0, 1) {
		return
	}
	defer atomic.StoreInt32(&m.mouseRun, 0)
	m.device.processMouseEvent(m.mouseEvent)
}
