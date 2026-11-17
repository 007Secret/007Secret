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

// authMessage is the first frame a peer sends so the extraction code stays out
// of the WebSocket URL (and therefore out of proxy / access logs).
type authMessage struct {
	Type     string `json:"type"`
	Key      string `json:"key"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type signalingRoom struct {
	sender        *websocket.Conn
	receiver      *websocket.Conn
	receiverReady bool
	writeMu       sync.Mutex
}

type SignalingHub struct {
	mu           sync.Mutex
	rooms        map[string]*signalingRoom
	upgrader     websocket.Upgrader
	pingInterval time.Duration
}

func NewSignalingHub() *SignalingHub {
	return &SignalingHub{
		rooms:        make(map[string]*signalingRoom),
		pingInterval: 20 * time.Second,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(*http.Request) bool { return true },
		},
	}
}

func (h *SignalingHub) Handle(c *gin.Context) {
	connection, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// Cap the frame size before any authentication so an unauthenticated peer
	// cannot exhaust server memory with an oversized frame. 64 KiB comfortably
	// fits auth and SDP signal frames.
	connection.SetReadLimit(64 * 1024)

	// The extraction code is never placed in the URL; the peer must authenticate
	// in its first frame so credentials do not appear in logs or browser history.
	_ = connection.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, authPayload, err := connection.ReadMessage()
	if err != nil {
		_ = connection.Close()
		return
	}
	var auth authMessage
	if err := json.Unmarshal(authPayload, &auth); err != nil ||
		auth.Type != "auth" || auth.Key == "" || auth.Password == "" ||
		(auth.Role != "sender" && auth.Role != "receiver") {
		_ = connection.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "invalid signaling auth"), time.Now().Add(time.Second))
		_ = connection.Close()
		return
	}

	roomKey := auth.Key + ":" + auth.Password
	role := auth.Role
	room, receiverReady, added := h.addPeer(roomKey, role, connection)
	if !added {
		_ = connection.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "role already connected"), time.Now().Add(time.Second))
		_ = connection.Close()
		return
	}
	// Browser WebSockets answer ping frames automatically. Keep an idle sender's
	// connection alive through reverse proxies while it waits for the recipient.
	const pongTimeout = 60 * time.Second
	_ = connection.SetReadDeadline(time.Now().Add(pongTimeout))
	connection.SetPongHandler(func(string) error {
		return connection.SetReadDeadline(time.Now().Add(pongTimeout))
	})
	heartbeatDone := make(chan struct{})
	go func() {
		ticker := time.NewTicker(h.pingInterval)
		defer ticker.Stop()
		for {
			select {
			case <-heartbeatDone:
				return
			case <-ticker.C:
				if err := connection.WriteControl(websocket.PingMessage, nil, time.Now().Add(5*time.Second)); err != nil {
					_ = connection.Close()
					return
				}
			}
		}
	}()
	defer close(heartbeatDone)
	if role == "sender" && receiverReady {
		h.write(room, connection, []byte(`{"type":"ready"}`))
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
		if err := json.Unmarshal(message, &signal); err != nil || !isSignalType(signal.Type) || (signal.Type != "ready" && len(signal.Payload) == 0) {
			continue
		}
		if signal.Type == "ready" && role != "receiver" {
			continue
		}
		h.forward(roomKey, role, signal.Type, message)
	}
}

func isSignalType(messageType string) bool {
	return messageType == "offer" || messageType == "answer" || messageType == "candidate" || messageType == "ready"
}

func (h *SignalingHub) addPeer(roomKey, role string, connection *websocket.Conn) (*signalingRoom, bool, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room := h.rooms[roomKey]
	if room == nil {
		room = &signalingRoom{}
		h.rooms[roomKey] = room
	}
	if role == "sender" {
		if room.sender != nil {
			return nil, false, false
		}
		room.sender = connection
		return room, room.receiverReady, true
	}
	if room.receiver != nil {
		return nil, false, false
	}
	room.receiver = connection
	return room, false, true
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

func (h *SignalingHub) forward(roomKey, role, signalType string, message []byte) {
	h.mu.Lock()
	room := h.rooms[roomKey]
	if room == nil {
		h.mu.Unlock()
		return
	}
	if signalType == "ready" {
		room.receiverReady = true
	}
	var peer *websocket.Conn
	if role == "sender" {
		peer = room.receiver
	} else {
		peer = room.sender
	}
	h.mu.Unlock()
	if peer != nil {
		h.write(room, peer, message)
	}
}

func (h *SignalingHub) write(room *signalingRoom, peer *websocket.Conn, message []byte) {
	room.writeMu.Lock()
	defer room.writeMu.Unlock()
	_ = peer.WriteMessage(websocket.TextMessage, message)
}
