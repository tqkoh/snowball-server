package streamer

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/tqkoh/snowball-server/streamer/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *streamer) ConnectWS(c echo.Context) error {
	roomID := c.QueryParam("room")
	if roomID == "" {
		roomID = "A"
	}

	connection, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}
	defer connection.Close()

	client := newClient(roomID, connection, s.receiver) // receiver is shared by streamer and all clients

	s.clients[client.id] = client
	go client.listen()
	go client.send()

	<-client.closer

	u, ok := users[client.id]
	name := "unknown"
	if ok {
		name = u.Name
	}
	processDead(s, client.id, client.id, fmt.Sprintf("%v disconnected", name), true)
	utils.Del(s.clients, client.id)

	return c.NoContent(http.StatusOK)
}
