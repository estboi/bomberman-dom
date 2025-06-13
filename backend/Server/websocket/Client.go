package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type ClientList map[*Client]bool

type Client struct {
	UserID     int
	Name       string
	Connection *websocket.Conn
	Manager    *Manager
	JoinTime   time.Time

	MessageChan chan Event
}

func NewClient(ws *websocket.Conn, m *Manager) *Client {
	return &Client{
		Connection:  ws,
		Manager:     m,
		MessageChan: make(chan Event),
	}
}

func (c *Client) ReadLoop() {
	defer func() {
		c.Manager.RemoveClient(c)
	}()

	if err := c.Connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("error setting timer: %s\n", err)
		return
	}

	c.Connection.SetPongHandler(c.pongHandler)

	c.Connection.SetReadLimit(10 * 1024 * 1024)

	for {
		_, payload, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				fmt.Printf("Websocket connection error: %s\n", err)
			} else {
				fmt.Printf("Websocket connection closed: %s\n", err)
			}
			return
		}
		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error unmarshaling event: %s\n", err)
			return
		}

		if err := c.Manager.routeEvent(request, c); err != nil {
			log.Panicf("error routing event: %s\n", err)
		}
	}
}

func (c *Client) WriteLoop() {
	defer func() {
		c.Manager.RemoveClient(c)
	}()

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case message, ok := <-c.MessageChan:
			if !ok {
				if err := c.Connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Printf("connection closed: %s\n", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("error marshaling data: %s\n", err)
			}

			if err := c.Connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("failed to send a message: %s\n", err)
			}
		case <-ticker.C:
			if err := c.Connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Printf("error pinging connection: %s\n", err)
				return
			}
		}
	}
}

func (c *Client) pongHandler(msg string) error {
	return c.Connection.SetReadDeadline(time.Now().Add(pongWait))
}
