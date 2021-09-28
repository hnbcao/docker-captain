package docker

import (
	"github.com/gorilla/websocket"
	"io"
)

type Writer struct {
	conn        *websocket.Conn
	messageType int
}

func (w Writer) Write(p []byte) (n int, err error) {
	n = len(p)
	err = nil
	if w.conn != nil {
		err = w.conn.WriteMessage(w.messageType, p)
		if err != nil {
			n = 0
		}
	}
	return
}

func NewWriter(conn *websocket.Conn, messageType int) io.Writer {
	return Writer{
		conn:        conn,
		messageType: messageType,
	}
}
