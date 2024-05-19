package websocket

import (
	"api/database"
	"api/database/models"
	"encoding/json"
	"errors"
	"log"

	uuid "github.com/satori/go.uuid"
)

type PlayerState struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	IsConnected bool      `json:"isConnected"`
	IsOwner     bool      `json:"isOwner"`
}

type GameState struct {
	Id             uuid.UUID     `json:"id"`
	IsRunning      bool          `json:"isRunning"`
	SongsToWin     uint          `json:"songsToWin"`
	TitleRequired  bool          `json:"titleRequired"`
	ArtistRequired bool          `json:"artistRequired"`
	Configured     bool          `json:"configured"`
	Players        []PlayerState `json:"players"`
}

func isConnected(playerId uuid.UUID, lobby map[*Client]bool) bool {
	for client := range lobby {
		if client.PlayerId == playerId {
			return true
		}
	}

	return false
}

func createPlayersUpdate(gameId uuid.UUID, pool *ConnectionPool) ([]PlayerState, error) {
	update := []PlayerState{}
	var game models.Game

	database := database.Get()
	playersError := database.Preload("Players").Where("id = ?", gameId).First(&game).Error
	if playersError != nil {
		return update, errors.New("could not find players")
	}

	lobby := pool.Lobbies[gameId]
	for _, player := range game.Players {
		update = append(update, PlayerState{
			Id:          player.Id,
			Name:        player.Name,
			IsConnected: isConnected(player.Id, lobby),
			IsOwner:     player.Id == game.OwnerId,
		})
	}

	return update, nil
}
func createGameUpdate(gameId uuid.UUID, pool *ConnectionPool) (GameState, error) {
	var game models.Game

	database := database.Get()
	gameError := database.Where("id = ?", gameId).First(&game).Error
	if gameError != nil {
		return GameState{}, errors.New("could not load game")
	}

	log.Println("GAME:", game.IsRunning)

	players, playersError := createPlayersUpdate(gameId, pool)
	if playersError != nil {
		return GameState{}, playersError
	}

	return GameState{
		Id:             game.Id,
		Configured:     game.Configured,
		IsRunning:      game.IsRunning,
		SongsToWin:     game.SongsToWin,
		TitleRequired:  game.TitleRequired,
		ArtistRequired: game.ArtistRequired,
		Players:        players,
	}, nil
}

func BroadcastGameUpdate(client *Client, pool *ConnectionPool) error {
	gameUpdate, updateError := createGameUpdate(client.GameId, pool)
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
