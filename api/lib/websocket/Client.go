package websocket

import (
	"time"

	gorilla "github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

const (
	maxWriteDuration = 10 * time.Second
	maxPongWait      = time.Minute
	pingInterval     = (maxPongWait * 9) / 10
	maxMessageSize   = 512
)

type Client struct {
	Connection       *gorilla.Conn
	GameId           uuid.UUID
	PlayerId         uuid.UUID
	OutgoingMessages chan ServerMessage
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

func (client *Client) Read(pool *ConnectionPool) {
	connection := client.Connection

	// Remove this client and unregister once we cannot continue reading messages.
	defer func() {
		pool.Remove <- client
		connection.Close()
	}()

	setConnectionSettings(connection)

	for {
		var message ClientMessage
		messageError := connection.ReadJSON(&message)

		if messageError != nil {
			connection.WriteJSON(ServerMessage{
				Type:    "error",
				Content: "Error reading message",
			})
			return
		}

		HandleClientMessage(message, client, pool)
	}
}

func (client *Client) Write(pool *ConnectionPool) {
	connection := client.Connection
	pingTimer := time.NewTicker(pingInterval)

	// Remove this client and unregister once we cannot continue writing messages.
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
		// Send pings to check the connection at each interval.
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
	client.OutgoingMessages <- message
}
