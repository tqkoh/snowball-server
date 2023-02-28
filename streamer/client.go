package streamer

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

type receiveData struct {
	id      uuid.UUID
	payload []byte
}

type client struct {
	id       uuid.UUID
	conn     *websocket.Conn
	receiver chan receiveData
	sender   chan string
	closer   chan bool
}

func newClient(roomID string, conn *websocket.Conn, receiver chan receiveData) *client {
	return &client{
		id:       uuid.Must(uuid.NewV4()),
		conn:     conn,
		receiver: receiver,
		sender:   make(chan string),
		closer:   make(chan bool),
	}
}
func (c *client) listen() {
	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			c.closer <- true
			return
		}
		if messageType != websocket.TextMessage {
			continue
		}
		fmt.Printf("message: %s\n", message)

		c.receiver <- receiveData{
			id:      c.id,
			payload: message,
		}
	}
}

func (c *client) send() {
	for {
		message := <-c.sender

		err := c.conn.WriteMessage(websocket.TextMessage, []byte(message))

		if err != nil {
			c.closer <- true
			return
		}
	}
}
