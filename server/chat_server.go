package server

import (
    "code.google.com/p/go.net/websocket"
    "errors"
)

type ChatServer struct {
    Room map[string][]*Client
}

type Client struct {
    Connection *websocket.Conn
}

var RoomNotExist = errors.New("Room not exist.")

func NewChatServer() *ChatServer {
    return &ChatServer{Room: make(map[string][]*Client)}
}

func (chatServer *ChatServer) CreateRoom(roomName string, creator *Client) {
    clients := []*Client{creator}

    chatServer.Room[roomName] = clients
}

func (chatServer *ChatServer) AddClient(c *Client, roomName string) (err error) {
    if room, exist := chatServer.Room[roomName]; exist {
        chatServer.Room[roomName] = append(room, c)
    } else {
        err = RoomNotExist
    }

    return
}

func (chatServer *ChatServer) FindClientFromRoom(connection *websocket.Conn, roomName string) (*Client, int) {
    for i, c := range chatServer.Room[roomName] {
        if c.Connection == connection {
            return c, i
        }
    }
    return nil, -1
}

func (chatServer *ChatServer) DeleteRoom(roomName string) {
    delete(chatServer.Room, roomName)
}

func NewClient(ws *websocket.Conn) *Client {
    return &Client{Connection: ws}
}
