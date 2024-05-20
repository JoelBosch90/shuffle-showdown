package websocket

import (
	"api/database"
	"api/database/models"
	"errors"

	uuid "github.com/satori/go.uuid"
)

type PlayerState struct {
	Id          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	IsConnected bool              `json:"isConnected"`
	IsOwner     bool              `json:"isOwner"`
	WonTracks   []models.WonTrack `json:"wonTracks"`
}

type GameState struct {
	Id             uuid.UUID      `json:"id"`
	IsRunning      bool           `json:"isRunning"`
	SongsToWin     uint           `json:"songsToWin"`
	TitleRequired  bool           `json:"titleRequired"`
	ArtistRequired bool           `json:"artistRequired"`
	Configured     bool           `json:"configured"`
	Players        []PlayerState  `json:"players"`
	Rounds         []models.Round `json:"rounds"`
}

func isConnected(playerId uuid.UUID, lobby map[*Client]bool) bool {
	for client := range lobby {
		if client.PlayerId == playerId {
			return true
		}
	}

	return false
}

func filterWonTracksByGameId(gameId uuid.UUID, wonTracks []models.WonTrack) []models.WonTrack {
	filtered := []models.WonTrack{}
	for _, wonTrack := range wonTracks {
		if wonTrack.GameId == gameId {
			filtered = append(filtered, wonTrack)
		}
	}

	return filtered
}

func createPlayersUpdate(gameId uuid.UUID, pool *ConnectionPool) ([]PlayerState, error) {
	update := []PlayerState{}
	var game models.Game

	database := database.Get()
	playersError := database.Preload("Players.WonTracks.Track.Artists").Where("id = ?", gameId).First(&game).Error
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
			WonTracks:   filterWonTracksByGameId(gameId, player.WonTracks),
		})
	}

	return update, nil
}

func lastRound(rounds []models.Round) (int, models.Round) {
	latestRound := models.Round{}
	latestIndex := 0

	for index, round := range rounds {
		if round.Number > latestRound.Number {
			latestRound = round
			latestIndex = index
		}
	}

	return latestIndex, latestRound
}

func hideTrackDetailsFromCurrentRound(rounds []models.Round) []models.Round {
	if len(rounds) == 0 {
		return rounds
	}

	currentRoundIndex, currentRound := lastRound(rounds)
	rounds[currentRoundIndex].Track = models.Track{
		PreviewUrl: currentRound.Track.PreviewUrl,
	}

	return rounds
}

func createGameUpdate(gameId uuid.UUID, pool *ConnectionPool) (GameState, error) {
	var game models.Game

	database := database.Get()
	gameError := database.Preload("Rounds.Track").Where("id = ?", gameId).First(&game).Error
	if gameError != nil {
		return GameState{}, errors.New("could not load game")
	}

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
		Rounds:         hideTrackDetailsFromCurrentRound(game.Rounds),
	}, nil
}

func BroadcastGameUpdate(client *Client, pool *ConnectionPool) error {
	gameUpdate, updateError := createGameUpdate(client.GameId, pool)
	if updateError != nil {
		return errors.New("invalid game state")
	}

	pool.Broadcast <- ServerMessage{
		Type:    ServerMessageTypeGameUpdate,
		Payload: gameUpdate,
		GameId:  client.GameId,
	}

	return nil
}
