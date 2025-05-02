package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var (
	HubInstance *Hub
	once        sync.Once
)

type Client struct {
	ID      string
	Conn    *websocket.Conn
	UserID  uuid.UUID
	Send    chan WsNotification
	AppType string
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan WsNotification
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

// WsNotification matches the Notification entity structure
type WsNotification struct {
	ID          uuid.UUID      `json:"id"`
	Application string         `json:"application"`
	Name        string         `json:"name"`
	URL         string         `json:"url"`
	ReadAt      *time.Time     `json:"read_at"`
	Message     string         `json:"message"`
	UserID      uuid.UUID      `json:"user_id"`
	CreatedBy   uuid.UUID      `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GetHub() *Hub {
	once.Do(func() {
		HubInstance = &Hub{
			broadcast:  make(chan WsNotification),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			clients:    make(map[*Client]bool),
		}
		go HubInstance.Run()
	})
	return HubInstance
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()

		case notification := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				if client.UserID == notification.UserID &&
					(notification.Application == "" || client.AppType == notification.Application) {
					select {
					case client.Send <- notification:
					default:
						log.Printf("error: client send channel is full")
						close(client.Send)
						delete(h.clients, client)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) BroadcastNotification(notification WsNotification) {
	h.broadcast <- notification
}

func (c *Client) readPump() {
	defer func() {
		GetHub().unregister <- c
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

func (c *Client) writePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case notification, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			json.NewEncoder(w).Encode(notification)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request, userID uuid.UUID, appType string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		ID:      uuid.New().String(),
		Conn:    conn,
		UserID:  userID,
		Send:    make(chan WsNotification, 256),
		AppType: appType,
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}
