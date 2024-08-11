package websocket

import uuid "github.com/satori/go.uuid"

type ServerMessageType string

const (
	ServerMessageTypeAnswerSelectionUpdate ServerMessageType = "answer-selection-update"
	ServerMessageTypeGameSessionUpdate     ServerMessageType = "game-session-update"
	ServerMessageTypeKickedPlayer          ServerMessageType = "kicked-player"
	ServerMessageTypeError                 ServerMessageType = "error"
	ServerMessageTypePlaylistLoadingUpdate ServerMessageType = "playlist-loading-update"
	ServerMessageTypePlaylistLoaded        ServerMessageType = "playlist-loaded"
)

type ErrorMessagePayload struct {
	Message string `json:"message"`
}

type ServerMessage struct {
	Type    ServerMessageType `json:"type"`
	Payload interface{}       `json:"payload"`
	GameId  uuid.UUID         `json:"-"`
}
