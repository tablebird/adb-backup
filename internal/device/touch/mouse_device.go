package touch

import (
	"adb-backup/internal/log"
	"adb-backup/internal/shell"
)

type mouseDevice interface {
	size() (int, int)
	processMouseEvent(mouseEvent chan *MouseEvent)
}

func newEventDevice(s shell.Shell) mouseDevice {
	exists := shell.WhichCommandExists(s, "su")
	if exists {
		return &shellEventDevice{s: s}
	} else {
		return &shellInputDevice{s: s}
	}
}

type shellInputDevice struct {
	s shell.Shell
}

func (t *shellInputDevice) size() (int, int) {
	x, y, err := shell.WmSize(t.s)
	if err != nil {
		return 0, 0
	}
	return x, y
}

func (t *shellInputDevice) processMouseEvent(mouseEvent chan *MouseEvent) {
	var downEvent *MouseEvent = nil
	for e := range mouseEvent {
		switch e.mouseType {
		case MouseTypeDown:
			downEvent = e
		case MouseTypeUp:
			if downEvent != nil {
				err := shell.InputSwipe(t.s, downEvent.x, downEvent.y, e.x, e.y, e.duration)
				if err != nil {
					log.ErrorF("inputSwipe err %s", err.Error())
				}
			}
			downEvent = nil
		default:
		}
	}
}

type shellEventDevice struct {
	s        shell.Shell
	touchDev string
}

func (t *shellEventDevice) size() (int, int) {
	dev, x, y, _ := shell.GetEventTouch(t.s)
	t.touchDev = dev
	return x, y
}

func (t *shellEventDevice) processMouseEvent(mouseEvent chan *MouseEvent) {
	for e := range mouseEvent {
		var builder = shell.NewSendEventBuilder(t.touchDev)
		appendMouseEvent(builder, e)
		length := len(mouseEvent)
		for i := 0; i < length; i++ {
			e := <-mouseEvent
			appendMouseEvent(builder, e)
		}
		if builder != nil {
			err := shell.SendEventOfBuilder(t.s, builder)
			if err != nil {
				log.WarningF("SendEvent error %v", err)
			}
		}
	}
}

func appendMouseEvent(builder shell.SendEventBuilder, e *MouseEvent) {
	switch e.mouseType {
	case MouseTypeDown:
		builder.Touch(e.id).X(e.x).Y(e.y).Pressure(1024)
	case MouseTypeMove:
		builder.X(int(e.x)).Y(int(e.y))
	case MouseTypeUp:
		builder.Pressure(0).ClearTouch()
	default:
		// Unknown action
		return
	}
	builder.Sync()
}
