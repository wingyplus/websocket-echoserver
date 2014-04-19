package server_test

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/wingyplus/websocket/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEchoServer(t *testing.T) {
	http.Handle("/echo", websocket.Handler(server.EchoServer))
	ts := httptest.NewServer(nil)
	defer ts.Close()

	serverAddr := ts.Listener.Addr().String()

	conn, err := websocket.Dial(
		fmt.Sprintf("ws://%s/echo", serverAddr),
		"tcp",
		fmt.Sprintf("http://%s", serverAddr))
	defer conn.Close()

	if err != nil {
		t.Errorf("Websocket err %s", err.Error())
	}

	msg := []byte("Hello World\n")
	if _, err := conn.Write(msg); err != nil {
		t.Errorf("Write to server error: %s", err.Error())
	}

	actualMessage := make([]byte, 512)
	n, err := conn.Read(actualMessage)
	if !bytes.Equal(msg, actualMessage[:n]) {
		t.Errorf("expect %v but was %v", msg, actualMessage)
	}
}
