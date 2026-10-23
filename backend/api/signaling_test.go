package api

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TestSignalingHubForwardsOfferOnlyToPairedReceiver(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	hub := NewSignalingHub()
	router.GET("/api/transfer/signal", hub.Handle)
	server := httptest.NewServer(router)
	defer server.Close()

	connect := func(role string) *websocket.Conn {
		u, err := url.Parse(server.URL)
		if err != nil {
			t.Fatalf("parse server URL: %v", err)
		}
		u.Scheme = "ws"
		u.Path = "/api/transfer/signal"
		u.RawQuery = url.Values{"key": {"secret"}, "password": {"pass"}, "role": {role}}.Encode()
		connection, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			t.Fatalf("connect %s: %v", role, err)
		}
		return connection
	}

	sender := connect("sender")
	defer sender.Close()
	receiver := connect("receiver")
	defer receiver.Close()

	message := `{"type":"offer","payload":{"sdp":"test-offer"}}`
	if err := sender.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("send offer: %v", err)
	}

	_, received, err := receiver.ReadMessage()
	if err != nil {
		t.Fatalf("read forwarded offer: %v", err)
	}
	if string(received) != message {
		t.Fatalf("received = %s, want %s", received, message)
	}
}

func TestSignalingHubDeliversStoredReceiverReadySignalToLaterSender(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	hub := NewSignalingHub()
	router.GET("/api/transfer/signal", hub.Handle)
	server := httptest.NewServer(router)
	defer server.Close()

	connect := func(role string) *websocket.Conn {
		endpoint := strings.Replace(server.URL, "http://", "ws://", 1) + "/api/transfer/signal?key=secret&password=pass&role=" + role
		connection, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
		if err != nil {
			t.Fatalf("connect %s: %v", role, err)
		}
		return connection
	}

	receiver := connect("receiver")
	defer receiver.Close()
	if err := receiver.WriteMessage(websocket.TextMessage, []byte(`{"type":"ready"}`)); err != nil {
		t.Fatalf("send ready: %v", err)
	}

	sender := connect("sender")
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
	gin.SetMode(gin.TestMode)
	router := gin.New()
	hub := NewSignalingHub()
	router.GET("/api/transfer/signal", hub.Handle)
	server := httptest.NewServer(router)
	defer server.Close()

	endpoint := strings.Replace(server.URL, "http://", "ws://", 1) + "/api/transfer/signal?key=secret&password=pass&role="
	if _, response, err := websocket.DefaultDialer.Dial(endpoint+"intruder", nil); err == nil || response == nil || response.StatusCode != 400 {
		t.Fatalf("unknown role must receive HTTP 400, response=%v err=%v", response, err)
	}

	sender, _, err := websocket.DefaultDialer.Dial(endpoint+"sender", nil)
	if err != nil {
		t.Fatalf("connect sender: %v", err)
	}
	defer sender.Close()

	thirdSender, _, err := websocket.DefaultDialer.Dial(endpoint+"sender", nil)
	if err != nil {
		t.Fatalf("connect third sender: %v", err)
	}
	defer thirdSender.Close()
	if _, _, err := thirdSender.ReadMessage(); !websocket.IsCloseError(err, websocket.ClosePolicyViolation) {
		t.Fatalf("third sender error = %v, want policy close", err)
	}
}

func TestSignalingHubKeepsWaitingSenderAlive(t *testing.T) {
	router := gin.New()
	hub := NewSignalingHub()
	hub.pingInterval = 10 * time.Millisecond
	router.GET("/api/transfer/signal", hub.Handle)
	server := httptest.NewServer(router)
	defer server.Close()

	endpoint := strings.Replace(server.URL, "http://", "ws://", 1) + "/api/transfer/signal?key=idle-test&password=test&role=sender"
	sender, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer sender.Close()

	pings := make(chan struct{}, 1)
	sender.SetPingHandler(func(data string) error {
		select {
		case pings <- struct{}{}:
		default:
		}
		return sender.WriteControl(websocket.PongMessage, []byte(data), time.Now().Add(time.Second))
	})
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
