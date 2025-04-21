package ws

import "github.com/kudras3r/KDTog/pkg/logger"

type Message struct {
	Sender  *Client
	Content []byte
}

type Hub struct {
	log *logger.Logger

	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
}

func NewHub(log *logger.Logger) *Hub {
	return &Hub{
		log:        log,
		broadcast:  make(chan Message, 1024),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	h.log.Info("hub started")
	for {
		select {
		case client := <-h.register:
			h.log.Infof("client %s registered", client.conn.RemoteAddr())
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.log.Infof("client %s unregistered", client.conn.RemoteAddr())
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				if client != message.Sender {
					h.log.Infof("sending message to client %s", client.conn.RemoteAddr())
					select {
					case client.send <- message.Content:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
