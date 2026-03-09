package socket

import (
	"errors"
	"io"
)

type Sender interface {
	io.Closer
}

type DateSender interface {
	Sender

	Send(data *data) (err error)
}

func NewSocketClientSender(client *socketClient) DateSender {
	return &socketClientSender{client: client}
}

type socketClientSender struct {
	client *socketClient
}

func (s *socketClientSender) Send(data *data) error {
	select {
	case s.client.send <- data:
		return nil
	default:
		return errors.New("send channel is closed or full")
	}
}

func (s *socketClientSender) Close() error {
	return s.client.Close()
}
