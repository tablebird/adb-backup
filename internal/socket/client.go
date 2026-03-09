package socket

import (
	"adb-backup/internal/log"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	pingWait = 60 * time.Second
)

type Client interface {
	io.Closer

	doWork()
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewClient(c *gin.Context) (Client, error) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}
	return &socketClient{
		conn: conn,
		send: make(chan *data, 10),
		read: make(chan *data, 100)}, nil
}

type data struct {
	messageType int
	data        []byte
}

type socketClient struct {
	conn *websocket.Conn
	send chan *data
	read chan *data
}

func (c *socketClient) doWork() {
	go c.writeLoop()
	go c.processRead()
	c.readLoop()
}

func (c *socketClient) Close() error {
	close(c.send)
	return c.conn.Close()
}

func (c *socketClient) readLoop() {

	c.setReadDeadline()
	c.conn.SetPingHandler(func(appData string) error {
		c.setReadDeadline()
		return nil
	})
	for {
		mt, message, err := c.conn.ReadMessage()
		if err != nil {
			log.WarningF("read socket error: %v", err)
			break
		}
		c.read <- &data{messageType: mt, data: message}
		c.setReadDeadline()
	}
}

func (c *socketClient) processRead() {
	sender := NewSocketClientSender(c)
	for d := range c.read {
		if d.messageType == websocket.TextMessage {
			processText(d.data, sender)
		}
	}
}

func (c *socketClient) setReadDeadline() {
	c.conn.SetReadDeadline(time.Now().Add(pingWait))
}

func (c *socketClient) writeLoop() {
	for message := range c.send {
		if err := c.conn.WriteMessage(message.messageType, message.data); err != nil {
			return
		}
	}

	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}
