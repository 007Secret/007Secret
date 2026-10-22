package api

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

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
