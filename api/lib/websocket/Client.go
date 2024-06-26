package websocket

import (
	"time"

	gorilla "github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

const (
	maxWriteDuration = 10 * time.Second
	maxPongWait      = 6 * maxWriteDuration           // Must be greater than write wait.
	pingInterval     = maxPongWait - maxWriteDuration // Must be less than pong wait.
	maxMessageSize   = 512
)

type Client struct {
	Pool             *ConnectionPool
	Connection       *gorilla.Conn
	GameId           uuid.UUID
	PlayerId         uuid.UUID
	OutgoingMessages chan ServerMessage
}

type ErrorMessagePayload struct {
	Message string `json:"message"`
}

func getReadDeadline() time.Time {
	return time.Now().Add(maxPongWait)
}

func getWriteDeadline() time.Time {
	return time.Now().Add(maxWriteDuration)
}

func setConnectionSettings(connection *gorilla.Conn) {
	connection.SetReadLimit(maxMessageSize)
	connection.SetReadDeadline(getReadDeadline())
	connection.SetPongHandler(func(string) error {
		connection.SetReadDeadline(getReadDeadline())
		return nil
	})
}

func (client *Client) Read() {
	connection := client.Connection

	// Close the client connection when we cannot continue reading messages.
	defer func() {
		// Remove the client from the connection pool.
		client.Pool.Remove <- client

		// Close the client connection.
		connection.Close()

		// Broadcast the updated game state.
		BroadcastGameUpdate(client, client.Pool)
	}()

	setConnectionSettings(connection)

	for {
		var message ClientMessage
		messageError := connection.ReadJSON(&message)

		if messageError != nil {
			connection.WriteJSON(ServerMessage{
				Type:    "error",
				Payload: "Error reading message",
			})
			return
		}

		HandleClientMessage(message, client, client.Pool)
	}
}

func (client *Client) Write() {
	connection := client.Connection
	pingTimer := time.NewTicker(pingInterval)

	// Stop the timer and close the client connection when we cannot continue writing messages.
	defer func() {
		pingTimer.Stop()
		connection.Close()
	}()

	for {
		select {
		// Send pending messages to the client.
		case message, ok := <-client.OutgoingMessages:
			connection.SetWriteDeadline(getWriteDeadline())
			if !ok {
				connection.WriteMessage(gorilla.CloseMessage, []byte{})
				return
			} else {
				writeError := connection.WriteJSON(message)
				if writeError != nil {
					break
				}
			}
		// Send pings to check the connection at each interval. This helps keep the connection alive.
		case <-pingTimer.C:
			connection.SetWriteDeadline(getWriteDeadline())
			pingError := connection.WriteMessage(gorilla.PingMessage, nil)
			if pingError != nil {
				return
			}
		}
	}
}

func (client *Client) Notify(message ServerMessage) {
	if client == nil {
		return
	}

	client.OutgoingMessages <- message
}

func (client *Client) SendError(message string) error {
	client.Notify(ServerMessage{
		Type:    ServerMessageTypeError,
		Payload: ErrorMessagePayload{Message: message},
	})

	return nil
}
