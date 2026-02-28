package main

import "github.com/gorilla/websocket"

// ClientList is a map used to help manage a map of clients
type ClientList map[*Client]bool 

//Client is a websocket client, basically a frontend visitor
type Client struct {
	//the websocket connection
	connection *websocket.Conn
	//manager is the manager used to manage the client
	manager *Manager
	// egress is used to avoid concurrent writes on the websocket
	egress chan []byte
}

// NewClient is used to initialize a new Client with all required values initialized
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager: manager,
		egress: make(chan []byte),
	}
}