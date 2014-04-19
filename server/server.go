package server

import (
	"code.google.com/p/go.net/websocket"
	"io"
)

func EchoServer(ws *websocket.Conn) {
	io.Copy(ws, ws)
}
