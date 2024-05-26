package websocket

import (
	"api/database"
	"api/database/models"
	gameHelpers "api/lib/game"
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
	Id          uuid.UUID      `json:"id"`
	HasStarted  bool           `json:"hasStarted"`
	HasFinished bool           `json:"hasFinished"`
	SongsToWin  uint           `json:"songsToWin"`
	Owner       models.Player  `json:"owner"`
	Players     []PlayerState  `json:"players"`
	Rounds      []models.Round `json:"rounds"`
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
			WonTracks:   gameHelpers.FilterWonTracksByGameId(gameId, player.WonTracks),
		})
	}

	return update, nil
}

func hideTrackDetailsFromCurrentRound(rounds []models.Round) []models.Round {
	if len(rounds) == 0 {
		return rounds
	}

	currentRoundIndex, currentRound := gameHelpers.LastRound(rounds)
	rounds[currentRoundIndex].Track = models.Track{
		PreviewUrl: currentRound.Track.PreviewUrl,
	}

	return rounds
}

func createGameUpdate(gameId uuid.UUID, pool *ConnectionPool) (GameState, error) {
	var game models.Game

	database := database.Get()
	gameError := database.Preload("Rounds.Track.Artists").Preload("Owner").Where("id = ?", gameId).First(&game).Error
	if gameError != nil {
		return GameState{}, errors.New("could not load game")
	}

	players, playersError := createPlayersUpdate(gameId, pool)
	if playersError != nil {
		return GameState{}, playersError
	}

	rounds := game.Rounds
	if !game.HasFinished {
		rounds = hideTrackDetailsFromCurrentRound(game.Rounds)
	}

	return GameState{
		Id:          game.Id,
		HasStarted:  game.HasStarted,
		HasFinished: game.HasFinished,
		SongsToWin:  game.SongsToWin,
		Owner:       game.Owner,
		Players:     players,
		Rounds:      rounds,
	}, nil
}

func BroadcastGameUpdate(client *Client, pool *ConnectionPool) error {
	gameUpdate, updateError := createGameUpdate(client.GameId, pool)
	if updateError != nil {
		return errors.New("invalid game state")
	}

	pool.Broadcast <- ServerMessage{
		Type:    ServerMessageTypeGameSessionUpdate,
		Payload: gameUpdate,
		GameId:  client.GameId,
	}

	return nil
}
