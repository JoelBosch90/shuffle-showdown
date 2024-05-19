package websocket

import (
	"api/database"
	"api/database/models"
	"encoding/json"
	"errors"

	uuid "github.com/satori/go.uuid"
)

type GameState struct {
	IsRunning      bool `json:"isRunning"`
	SongsToWin     uint `json:"songsToWin"`
	TitleRequired  bool `json:"titleRequired"`
	ArtistRequired bool `json:"artistRequired"`
}

func createGameUpdate(gameId uuid.UUID) (GameState, error) {
	var game models.Game

	database := database.Get()
	gameError := database.Where("id = ?", gameId).First(&game).Error
	if gameError != nil {
		return GameState{}, errors.New("could not load game")
	}

	return GameState{
		IsRunning:      game.IsRunning,
		SongsToWin:     game.SongsToWin,
		TitleRequired:  game.TitleRequired,
		ArtistRequired: game.ArtistRequired,
	}, nil
}

func BroadcastGameUpdate(client *Client, pool *ConnectionPool) error {
	gameUpdate, updateError := createGameUpdate(client.GameId)
	if updateError != nil {
		return errors.New("invalid game state")
	}

	updateJson, jsonError := json.Marshal(&gameUpdate)
	if jsonError != nil {
		return errors.New("invalid game state")
	}

	pool.Broadcast <- ServerMessage{
		Type:    ServerMessageTypeGameUpdate,
		Payload: string(updateJson),
		GameId:  client.GameId,
	}

	return nil
}
