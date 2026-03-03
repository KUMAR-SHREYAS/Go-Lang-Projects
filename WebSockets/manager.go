package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	/**
	websocketUpgrader is used to upgrade incomming HTTP requests into a persitent websocket connection
	*/
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

var (
	ErrEventNotSupported  = errors.New("this event type is not supported")	
)

// Manager is used to hold references to all Clients Registered, and Broadcasting etc
type Manager struct {
	clients ClientList
	// Using a syncMutex here to be able to lcok state before editing clients
	// Could also use Channels to block
	sync.RWMutex
	// handlers are functions that are used to handle Events
	handlers map[string]EventHandler

}

// NewManager is used to initalize all the values inside the manager
func NewManager() *Manager {
	m:= &Manager{
		clients: make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

// setupEventHandlers configures and adds all handlers
func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = func (e Event, c *Client) error {
		fmt.Println(e)
		return nil
	}
}

// routeEvent is used to make sure the correct event goes into the correct handler
func (m *Manager) routeEvent(event Event, c *Client) error {
	// Check if Handler is present in Map
	if handler, ok := m.handlers[event.Type]; ok {
		// Execute the handler and return any err
		if err := handler(event, c)	; err != nil {
			return err
		}
		return nil
	}else {
		return ErrEventNotSupported
	}
}

// serveWS is a HTTP Handler that the has the Manager that allows connections
func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")
	// Begin by upgrading the HTTP request
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Create New Client
	client := NewClient(conn, m)
	//Add the newly created client to the manager
	m.addClient(client)
	// The WebSocket connection is only allowed to have one concurrent writer,
	//  we can fix this by having an unbuffered channel act as a locker.

	//Start the read/Write process
	go client.readMessages()
	go client.writeMessages()
}

// addClient will add clients to our clientList
func (m *Manager) addClient(client *Client) {
	// Lock so we can manipulate
	m.Lock()
	defer m.Unlock()

	// Add Client
	m.clients[client] = true
}

// removeClient will remove the client and clean up
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {

		// Remove from manager FIRST
		delete(m.clients, client)

		//  Stop write goroutine
		close(client.egress)

		// Close websocket connection
		client.connection.Close()
	}
}

// // readMessages will start the client to read messages and handle them
// // appropriately.
// // This is suppose to be ran as a goroutine
// func (c *Client) readMessages() {
// 	defer func() {
// 		// Graceful Close the Connection once this
// 		// function is done
// 		c.manager.removeClient(c)
// 	}()

// 	//Loop forevr
// 	for {
// 		// ReadMessage is used to read the next message in queue
// 		// in the connection
// 		messageType, payload, err := c.connection.ReadMessage()

// 		if err != nil {
// 			// If Connection is closed, we will Recieve an error here
// 			// We only want to log Strange errors, but not simple Disconnection
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("error reading message: %v", err)

// 				break
// 			}
// 		}
// 		log.Println("MessageType: ", messageType)
// 		log.Println("Payload:  ", string(payload))

// 		// Hack to test that WriteMessages works as intended
// 		// Will be replaced soon
// 		c.manager.RLock()
// 		for wsclient := range c.manager.clients {
// 			wsclient.egress <- payload
// 		}
// 		c.manager.RUnlock()
// 	}
// }

// // writeMessages is a process that listens for new messages to output to the Client
// func (c *Client) writeMessages() {
// 	defer func() {
// 		// Graceful close if this triggers a closing
// 		c.manager.removeClient(c)
// 	}()

// 	for {
// 		select {
// 		case message, ok := <-c.egress:
// 			// / Ok will be false Incase the egress channel is closed
// 			if !ok {
// 				// Manager has closed this connection channel, so communicate that to frontend
// 				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
// 					//Log that the connection is closed and the reason
// 					log.Println("connection closed: ", err)
// 				}
// 				// Return to close the goroutine
// 				return
// 			}
// 			//Write a Regular text message to the connection
// 			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
// 				log.Println(err)
// 			}
// 			log.Println("sent message")
// 		}
// 	}

// }
