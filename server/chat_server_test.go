package server_test

import (
	"code.google.com/p/go.net/websocket"
	"github.com/wingyplus/websocket/server"
	"testing"
)

func TestChatServerCreateRoom(t *testing.T) {
	var chatServer *server.ChatServer = server.NewChatServer()

	chatServer.CreateRoom("Hello", new(server.Client))

	if _, ok := chatServer.Room["Hello"]; !ok {
		t.Error("expect room name Hello but was nil")
	}
}

func TestChatServerAddClientToExistingRoom(t *testing.T) {
	var chatServer *server.ChatServer = server.NewChatServer()
	chatServer.CreateRoom("Hello", new(server.Client))

	chatServer.AddClient(new(server.Client), "Hello")

	v, ok := chatServer.Room["Hello"]

	if !ok {
		t.Error("expect room name Hello but was nil")
	}

	if len(v) != 2 {
		t.Error("client not add")
	}
}

func TestChatServerAddClientToNotExistingRoom(t *testing.T) {
	var chatServer *server.ChatServer = server.NewChatServer()

	err := chatServer.AddClient(new(server.Client), "Hello")

	if err == nil {
		t.Error("error not raise.")
	} else {
		if err.Error() != "Room not exist." {
			t.Error("room must be not exist.")
		}
	}
}

func TestDeleteRoom(t *testing.T) {
	var chatServer *server.ChatServer = server.NewChatServer()

	chatServer.CreateRoom("Hello", new(server.Client))

	chatServer.DeleteRoom("Hello")

	if _, ok := chatServer.Room["Hello"]; ok {
		t.Error("expect room name Hello not exist but exist")
	}
}

func TestClient(t *testing.T) {
	ws := new(websocket.Conn)
	var c *server.Client = server.NewClient(ws)

	if c.Connection != ws {
		t.Errorf("connection address %p is not same ws %p", c.Connection, ws)
	}
}

func TestFindClient(t *testing.T) {
	connections := []*websocket.Conn{}

	for i := 1; i <= 10; i++ {
		connections = append(connections, new(websocket.Conn))
	}

	chatServer := server.NewChatServer()
	chatServer.CreateRoom("Hello", server.NewClient(connections[0]))

	for i := 1; i < 10; i++ {
		chatServer.AddClient(server.NewClient(connections[i]), "Hello")
	}

	client, index := chatServer.FindClientFromRoom(connections[2], "Hello")

	if client.Connection != connections[2] {
		t.Errorf("connection address %p is not same ws %p", client.Connection, connections[2])
	}
	if index != 2 {
		t.Errorf("expect 2 but was %d", index)
	}

	client, index = chatServer.FindClientFromRoom(connections[4], "Hello")

	if client.Connection != connections[4] {
		t.Errorf("connection address %p is not same ws %p", client.Connection, connections[4])
	}
	if index != 4 {
		t.Errorf("expect 4 but was %d", index)
	}

	client, index = chatServer.FindClientFromRoom(new(websocket.Conn), "Hello")

	if client != nil {
		t.Errorf("expect nil but was %v", client)
	}
	if index != -1 {
		t.Errorf("expect -1 but was %d", index)
	}
}

func TestRemoveClientFromRoom(t *testing.T) {
	connections := []*websocket.Conn{}

	for i := 1; i <= 10; i++ {
		connections = append(connections, new(websocket.Conn))
	}

	chatServer := server.NewChatServer()
	chatServer.CreateRoom("Hello", server.NewClient(connections[0]))

	for i := 1; i < 10; i++ {
		chatServer.AddClient(server.NewClient(connections[i]), "Hello")
	}

	chatServer.RemoveClientFromRoom(connections[4], "Hello")

	if len(chatServer.Room["Hello"]) != 9 {
		t.Errorf("expect 9 but was %d", len(chatServer.Room["Hello"]))
	}

	c, i := chatServer.FindClientFromRoom(connections[4], "Hello")

	if c != nil {
		t.Errorf("expect nil but was %v", c)
	}
	if i != -1 {
		t.Errorf("expect -1 but was %d", i)
	}
}
