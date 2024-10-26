package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Struct to represent connected clients
type Client struct {
	conn *websocket.Conn
}

// Slice to hold connected clients (use a mutex for thread safety)
var clients = make([]*Client, 0)
var mutex = &sync.Mutex{} // Ensures safe concurrent access to the clients slice

// Function to broadcast a message to all clients except the sender
func broadcastMessage(sender *Client, msg []byte) {
	mutex.Lock() // Lock the clients list during broadcast
	defer mutex.Unlock()

	// Loop through all clients and send message to everyone except the sender
	for _, client := range clients {
		if client != sender {
			if err := client.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				fmt.Println("Error sending message:", err)
			}
		}
	}
}

// WebSocket handler
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // Upgrade HTTP to WebSocket
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// Create a new client
	client := &Client{conn: conn}

	// Add the new client to the list
	mutex.Lock()
	clients = append(clients, client)
	mutex.Unlock()

	// Listen for incoming messages from this client
	for {
		messageType, msg, err := conn.ReadMessage() // Read incoming message
		if err != nil {
			// Handle client disconnect (remove from client list)
			fmt.Println("Client disconnected")
			mutex.Lock()
			// Remove the client from the clients list
			for i, c := range clients {
				if c == client {
					clients = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			mutex.Unlock()
			return
		}

		// Broadcast the message to all other clients
		broadcastMessage(client, msg)

		// Optionally: Echo the message back to the sender too (if you want)
		if err := conn.WriteMessage(messageType, msg); err != nil {
			fmt.Println("Error echoing message to sender:", err)
			return
		}
	}
}

func main() {
	// Set up WebSocket route
	http.HandleFunc("/connect", wsHandler)

	// Start server on port 8080
	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}
