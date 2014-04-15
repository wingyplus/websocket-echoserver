package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io"
	"net/http"
)

type Server struct {
	Clients map[*websocket.Conn]string
	Room    map[string][]*websocket.Conn
}

func (s *Server) Add(roomName string, c *websocket.Conn) {
	s.Room[roomName] = append(s.Room[roomName], c)
	s.Clients[c] = roomName
}

func (s *Server) Remove(ws *websocket.Conn) {
	roomName := s.Clients[ws]
	delete(s.Clients, ws)

	clients := s.Room[roomName]
	for i, c := range clients {
		if ws == c {
			s.Room[roomName] = append(s.Room[roomName][:i], s.Room[roomName][i+1:]...)
		}
	}
	if len(s.Room[roomName]) == 0 {
		delete(s.Room, roomName)
	}
}

var server = &Server{
	Clients: make(map[*websocket.Conn]string),
	Room:    make(map[string][]*websocket.Conn),
}

func EchoServer(ws *websocket.Conn) {
	var req map[string]string

    defer server.Remove(ws)

	for websocket.JSON.Receive(ws, &req) != io.EOF {
		fmt.Println(ws.Request().RemoteAddr)
		switch req["event"] {
		case "ADD":
			roomName := req["roomName"]
			server.Add(roomName, ws)
			fmt.Println(server.Room)
		case "ECHO":
			for _, client := range server.Room[req["roomName"]] {
				websocket.Message.Send(client, req["message"])
			}
		case "CLOSE":
			fmt.Println(req["message"])
		}
	}

	fmt.Println("connection ", ws)
	fmt.Println(server.Clients)
	fmt.Println(server.Room)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("www")))
	http.Handle("/echo", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
