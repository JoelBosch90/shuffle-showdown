package websocket

import (
	uuid "github.com/satori/go.uuid"
)

type ConnectionPool struct {
	// This map stores a map of clients for each game identifier.
	// In this map of clients, each client's memory address points to a boolean.
	// The boolean is irrelevant, but using the map allows us easy registration
	// of client memory addresses.
	Lobbies   map[uuid.UUID]map[*Client]struct{}
	Remove    chan *Client
	Register  chan *Client
	Broadcast chan ServerMessage
	Notify    chan ServerMessage
}

var connectionPool = createConnectionPool()

func createConnectionPool() *ConnectionPool {
	pool := &ConnectionPool{
		Lobbies:   make(map[uuid.UUID]map[*Client]struct{}),
		Remove:    make(chan *Client),
		Register:  make(chan *Client),
		Broadcast: make(chan ServerMessage),
		Notify:    make(chan ServerMessage),
	}

	go pool.Run()

	return pool
}

// Function to expose the connection pool. We generally only want one.
func GetConnectionPool() *ConnectionPool {
	return connectionPool
}

func (pool *ConnectionPool) Run() {
	for {
		select {
		case connection := <-pool.Register:
			pool.RegisterConnection(connection)
		case client := <-pool.Remove:
			pool.RemoveConnection(client)
		case message := <-pool.Broadcast:
			pool.BroadcastMessage(message)
		}
	}
}

func (pool *ConnectionPool) RegisterConnection(client *Client) {
	// Create a new lobby if we didn't have one yet.
	if pool.Lobbies[client.GameId] == nil {
		pool.Lobbies[client.GameId] = make(map[*Client]struct{})
	}

	// This is only symbolic, to signal that a client is connected.
	pool.Lobbies[client.GameId][client] = struct{}{}
}

func (pool *ConnectionPool) RemoveConnection(client *Client) {
	gameLobby := pool.Lobbies[client.GameId]

	if gameLobby != nil {
		delete(gameLobby, client)
		close(client.OutgoingMessages)
	}
}

func (pool *ConnectionPool) BroadcastMessage(message ServerMessage) {
	gameLobby := pool.Lobbies[message.GameId]

	for client := range gameLobby {
		select {
		// Send the same message to all clients in the lobby.
		case client.OutgoingMessages <- message:
		// Remove the connection from the pool if we cannot send
		// the message.
		default:
			pool.RemoveConnection(client)
		}
	}
}
