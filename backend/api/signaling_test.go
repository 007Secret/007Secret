package api

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func newSignalingTestServer(t *testing.T) (*httptest.Server, *SignalingHub) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	hub := NewSignalingHub()
	router.GET("/api/transfer/signal", hub.Handle)
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)
	return server, hub
}

// dialRaw opens the WebSocket without any credentials in the URL.
func dialRaw(t *testing.T, serverURL string) *websocket.Conn {
	t.Helper()
	endpoint := strings.Replace(serverURL, "http://", "ws://", 1) + "/api/transfer/signal"
	connection, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		t.Fatalf("dial signaling: %v", err)
	}
	return connection
}

func sendAuth(t *testing.T, connection *websocket.Conn, key, password, role string) {
	t.Helper()
	auth, err := json.Marshal(map[string]string{"type": "auth", "key": key, "password": password, "role": role})
	if err != nil {
		t.Fatalf("marshal auth: %v", err)
	}
	if err := connection.WriteMessage(websocket.TextMessage, auth); err != nil {
		t.Fatalf("send auth: %v", err)
	}
}

// dialSignaling opens the socket and performs the first-frame auth handshake.
func dialSignaling(t *testing.T, serverURL, key, password, role string) *websocket.Conn {
	t.Helper()
	connection := dialRaw(t, serverURL)
	sendAuth(t, connection, key, password, role)
	return connection
}

func TestSignalingHubForwardsOfferOnlyToPairedReceiver(t *testing.T) {
	server, _ := newSignalingTestServer(t)

	sender := dialSignaling(t, server.URL, "secret", "pass", "sender")
	defer sender.Close()
	receiver := dialSignaling(t, server.URL, "secret", "pass", "receiver")
	defer receiver.Close()

	// The receiver announces readiness; the sender only sends its offer after
	// seeing this, which guarantees the receiver has joined the room.
	if err := receiver.WriteMessage(websocket.TextMessage, []byte(`{"type":"ready"}`)); err != nil {
		t.Fatalf("send ready: %v", err)
	}
	_ = sender.SetReadDeadline(time.Now().Add(time.Second))
	if _, ready, err := sender.ReadMessage(); err != nil || string(ready) != `{"type":"ready"}` {
		t.Fatalf("sender ready = %s, %v; want ready", ready, err)
	}

	message := `{"type":"offer","payload":{"sdp":"test-offer"}}`
	if err := sender.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("send offer: %v", err)
	}

	_ = receiver.SetReadDeadline(time.Now().Add(time.Second))
	_, received, err := receiver.ReadMessage()
	if err != nil {
		t.Fatalf("read forwarded offer: %v", err)
	}
	if string(received) != message {
		t.Fatalf("received = %s, want %s", received, message)
	}
}

func TestSignalingHubDeliversStoredReceiverReadySignalToLaterSender(t *testing.T) {
	server, _ := newSignalingTestServer(t)

	receiver := dialSignaling(t, server.URL, "secret", "pass", "receiver")
	defer receiver.Close()
	if err := receiver.WriteMessage(websocket.TextMessage, []byte(`{"type":"ready"}`)); err != nil {
		t.Fatalf("send ready: %v", err)
	}
	// Give the hub a moment to record receiverReady before the sender connects.
	time.Sleep(100 * time.Millisecond)

	sender := dialSignaling(t, server.URL, "secret", "pass", "sender")
	defer sender.Close()
	_ = sender.SetReadDeadline(time.Now().Add(time.Second))
	_, received, err := sender.ReadMessage()
	if err != nil {
		t.Fatalf("read stored ready signal: %v", err)
	}
	if string(received) != `{"type":"ready"}` {
		t.Fatalf("received = %s, want ready signal", received)
	}
}

func TestSignalingHubRejectsUnknownRoleAndThirdPeer(t *testing.T) {
	server, _ := newSignalingTestServer(t)

	// Unknown role: the socket upgrades, then auth is rejected and closed.
	intruder := dialRaw(t, server.URL)
	sendAuth(t, intruder, "secret", "pass", "intruder")
	_ = intruder.SetReadDeadline(time.Now().Add(time.Second))
	if _, _, err := intruder.ReadMessage(); !websocket.IsCloseError(err, websocket.ClosePolicyViolation) {
		t.Fatalf("unknown role error = %v, want policy close", err)
	}
	_ = intruder.Close()

	sender := dialSignaling(t, server.URL, "secret", "pass", "sender")
	defer sender.Close()
	time.Sleep(100 * time.Millisecond)

	thirdSender := dialSignaling(t, server.URL, "secret", "pass", "sender")
	defer thirdSender.Close()
	_ = thirdSender.SetReadDeadline(time.Now().Add(time.Second))
	if _, _, err := thirdSender.ReadMessage(); !websocket.IsCloseError(err, websocket.ClosePolicyViolation) {
		t.Fatalf("third sender error = %v, want policy close", err)
	}
}

func TestSignalingHubKeepsWaitingSenderAlive(t *testing.T) {
	server, hub := newSignalingTestServer(t)
	hub.pingInterval = 10 * time.Millisecond

	sender := dialRaw(t, server.URL)
	defer sender.Close()

	pings := make(chan struct{}, 1)
	sender.SetPingHandler(func(data string) error {
		select {
		case pings <- struct{}{}:
		default:
		}
		return sender.WriteControl(websocket.PongMessage, []byte(data), time.Now().Add(time.Second))
	})
	// Authenticate only after installing the ping handler so no heartbeat is missed.
	sendAuth(t, sender, "idle-test", "test", "sender")

	readDone := make(chan struct{})
	go func() {
		defer close(readDone)
		_, _, _ = sender.ReadMessage()
	}()
	defer func() { _ = sender.Close(); <-readDone }()

	select {
	case <-pings:
	case <-time.After(time.Second):
		t.Fatal("waiting sender received no heartbeat; an idle reverse proxy will close it")
	}
}

// TestSignalingHubRejectsOversizedUnauthenticatedFrame guards the pre-auth
// read limit: an oversized frame must terminate the connection instead of
// being buffered into memory.
func TestSignalingHubRejectsOversizedUnauthenticatedFrame(t *testing.T) {
	server, _ := newSignalingTestServer(t)

	intruder := dialRaw(t, server.URL)
	defer intruder.Close()

	oversized := make([]byte, 65*1024)
	if err := intruder.WriteMessage(websocket.TextMessage, oversized); err != nil {
		t.Fatalf("send oversized frame: %v", err)
	}
	_ = intruder.SetReadDeadline(time.Now().Add(2 * time.Second))
	if _, _, err := intruder.ReadMessage(); !websocket.IsCloseError(err, websocket.CloseMessageTooBig) {
		t.Fatalf("oversized frame error = %v, want close 1009 (message too big)", err)
	}
}
