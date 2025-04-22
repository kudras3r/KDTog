package ws

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kudras3r/KDTog/pkg/config"
	"github.com/kudras3r/KDTog/pkg/logger"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var maxMessageSize int64

var log *logger.Logger

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader websocket.Upgrader

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

type OutgoingMessage struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

func SetLogger(l *logger.Logger) {
	log = l
}

func SetConfig(config config.WSock) {
	upgrader = websocket.Upgrader{
		ReadBufferSize:  config.RWBuffSize,
		WriteBufferSize: config.RWBuffSize,
	}
	maxMessageSize = int64(config.MaxMessSize)
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		log.Infof("client %s disconnected", c.conn.RemoteAddr())
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		log.Debugf("pong received from client %s", c.conn.RemoteAddr())
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("unexpected close error from client %s: %v", c.conn.RemoteAddr(), err)
			} else {
				log.Infof("client %s disconnected: %v", c.conn.RemoteAddr(), err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Infof("message received from client %s: %s", c.conn.RemoteAddr(), string(message))
		c.hub.broadcast <- Message{Sender: c, Content: message}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		log.Infof("stopped writePump for client %s", c.conn.RemoteAddr())
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				log.Infof("hub closed the send channel for client %s", c.conn.RemoteAddr())
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			outgoingMessage := OutgoingMessage{
				Sender:  c.conn.RemoteAddr().String(),
				Content: string(message),
			}
			jsonMessage, err := json.Marshal(outgoingMessage)
			if err != nil {
				log.Errorf("error marshalling message for client %s: %v", c.conn.RemoteAddr(), err)
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Errorf("error creating writer for client %s: %v", c.conn.RemoteAddr(), err)
				return
			}
			w.Write(jsonMessage)
			if err := w.Close(); err != nil {
				log.Errorf("error closing writer for client %s: %v", c.conn.RemoteAddr(), err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Errorf("error sending ping to client %s: %v", c.conn.RemoteAddr(), err)
				return
			}
			log.Debugf("ping sent to client %s", c.conn.RemoteAddr())
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("failed to upgrade connection: %v", err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	log.Infof("new client connected: %s", conn.RemoteAddr())

	go client.writePump()
	go client.readPump()
}
