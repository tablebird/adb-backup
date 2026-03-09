package shell

import (
	"fmt"
	"strings"
)

const (
	_SENDEVENT = "sendevent"
)

type SendEventBuilder interface {
	Builder
	Event(t int, code int, value int) SendEventBuilder
	X(x int) SendEventBuilder
	Y(x int) SendEventBuilder
	Pressure(p int) SendEventBuilder
	Touch(touchId int) SendEventBuilder
	ClearTouch() SendEventBuilder
	Down() SendEventBuilder
	Up() SendEventBuilder
	Sync() SendEventBuilder
}

func NewSendEventBuilder(device string) SendEventBuilder {
	return &realSendEventBuilder{device: device}
}

type realSendEventBuilder struct {
	cmdBuilder
	device string
}

func (b *realSendEventBuilder) Event(t int, code int, value int) SendEventBuilder {
	b.multiAppend(buildSendEventCommand(b.device, t, code, value))
	return b
}

func (b *realSendEventBuilder) Touch(touchId int) SendEventBuilder {
	return b.Event(3, 57, touchId)
}

func (b *realSendEventBuilder) ClearTouch() SendEventBuilder {
	return b.Event(3, 57, 4294967295)
}

func (b *realSendEventBuilder) X(x int) SendEventBuilder {
	return b.Event(3, 53, x)
}

func (b *realSendEventBuilder) Y(y int) SendEventBuilder {
	return b.Event(3, 54, y)
}

func (b *realSendEventBuilder) Pressure(p int) SendEventBuilder {
	return b.Event(3, 58, p)
}

func (b *realSendEventBuilder) Down() SendEventBuilder {
	return b.Event(1, 330, 1)
}

func (b *realSendEventBuilder) Up() SendEventBuilder {
	return b.Event(1, 330, 0)
}

func (b *realSendEventBuilder) Sync() SendEventBuilder {
	return b.Event(0, 0, 0)
}

func buildSendEventCommand(dev string, t int, code int, value int) string {
	return fmt.Sprintf("su input %s %s %d %d %d", _SENDEVENT, dev, t, code, value)
}

func SendEvent(s Shell, dev string, t int, code int, value int) error {
	res, err := s.RunCommand(buildSendEventCommand(dev, t, code, value))
	if err != nil {
		return err
	}
	if strings.TrimSpace(res) != "" {
		return fmt.Errorf("sendEvent error: %s", res)
	}
	return nil
}

func SendEventOfBuilder(s Shell, builder SendEventBuilder) error {
	if asyncShell, ok := s.(AsyncShell); ok {
		return asyncShell.AsyncRunCommand(builder.Build())
	}
	res, err := RunBuilder(s, builder)
	if err != nil {
		return err
	}
	if strings.TrimSpace(res) != "" {
		return fmt.Errorf("sendEvent error: %s", res)
	}
	return nil
}
