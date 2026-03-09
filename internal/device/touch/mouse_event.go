package touch

type MouseEvent struct {
	id        int
	mouseType MouseType
	x, y      int
	duration  int
}

func NewMouseEvent(id int, mouseType MouseType, x, y int) *MouseEvent {
	return &MouseEvent{
		id:        id,
		mouseType: mouseType,
		x:         x,
		y:         y,
		duration:  100,
	}
}

func (m *MouseEvent) SetDuration(duration int) {
	m.duration = duration
}
