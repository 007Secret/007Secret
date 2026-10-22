package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type signalingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type signalingRoom struct {
	sender   *websocket.Conn
	receiver *websocket.Conn
}

type SignalingHub struct {
	mu       sync.Mutex
	rooms    map[string]*signalingRoom
	upgrader websocket.Upgrader
}

func NewSignalingHub() *SignalingHub {
	return &SignalingHub{
		rooms: make(map[string]*signalingRoom),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(*http.Request) bool { return true },
		},
	}
}

func (h *SignalingHub) Handle(c *gin.Context) {
	key := c.Query("key")
	password := c.Query("password")
	role := c.Query("role")
	if key == "" || password == "" || (role != "sender" && role != "receiver") {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid signaling request"})
		return
	}

	connection, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	roomKey := key + ":" + password
	if !h.addPeer(roomKey, role, connection) {
		_ = connection.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "role already connected"), time.Now().Add(time.Second))
		_ = connection.Close()
		return
	}
	defer func() {
		h.removePeer(roomKey, role, connection)
		_ = connection.Close()
	}()

	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			return
		}
		if messageType != websocket.TextMessage {
			continue
		}

		var signal signalingMessage
		if err := json.Unmarshal(message, &signal); err != nil || !isSignalType(signal.Type) || len(signal.Payload) == 0 {
			continue
		}
		if peer := h.peer(roomKey, role); peer != nil {
			_ = peer.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func isSignalType(messageType string) bool {
	return messageType == "offer" || messageType == "answer" || messageType == "candidate"
}

func (h *SignalingHub) addPeer(roomKey, role string, connection *websocket.Conn) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	room := h.rooms[roomKey]
	if room == nil {
		room = &signalingRoom{}
		h.rooms[roomKey] = room
	}
	if role == "sender" {
		if room.sender != nil {
			return false
		}
		room.sender = connection
		return true
	}
	if room.receiver != nil {
		return false
	}
	room.receiver = connection
	return true
}

func (h *SignalingHub) removePeer(roomKey, role string, connection *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room := h.rooms[roomKey]
	if room == nil {
		return
	}
	if role == "sender" && room.sender == connection {
		room.sender = nil
	}
	if role == "receiver" && room.receiver == connection {
		room.receiver = nil
	}
	if room.sender == nil && room.receiver == nil {
		delete(h.rooms, roomKey)
	}
}

func (h *SignalingHub) peer(roomKey, role string) *websocket.Conn {
	h.mu.Lock()
	defer h.mu.Unlock()

	room := h.rooms[roomKey]
	if room == nil {
		return nil
	}
	if role == "sender" {
		return room.receiver
	}
	return room.sender
}
