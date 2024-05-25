package websocket

import (
	"api/database"
	"api/database/models"
	gameHelpers "api/lib/game"
	"encoding/json"
	"errors"
)

type Answer struct {
	BeforeReleaseYear *int `json:"beforeReleaseYear"`
	AfterReleaseYear  *int `json:"afterReleaseYear"`
}

func SubmitAnswerHandler(message ClientMessage, client *Client, pool *ConnectionPool) error {
	database := database.Get()
	var game models.Game

	gameError := database.Preload("Rounds.Track").Where("id = ?", client.GameId).First(&game).Error
	if gameError != nil {
		return errors.New("could not find game")
	}

	_, currentRound := gameHelpers.LastRound(game.Rounds)
	if currentRound.PlayerId != client.PlayerId {
		return errors.New("not your turn")
	}

	var answer Answer
	answerParseError := json.Unmarshal([]byte(message.Payload), &answer)
	if answerParseError != nil {
		return errors.New("could not parse answer")
	}
	if answer.BeforeReleaseYear == nil && answer.AfterReleaseYear == nil {
		return errors.New("answer must contain at least one field")
	}

	correctBefore := answer.BeforeReleaseYear == nil || currentRound.Track.ReleaseYear <= uint(*answer.BeforeReleaseYear)
	correctAfter := answer.AfterReleaseYear == nil || currentRound.Track.ReleaseYear >= uint(*answer.AfterReleaseYear)
	if correctBefore && correctAfter {
		awardError := gameHelpers.AwardTrack(game.Id, currentRound.Track, models.Player{Id: client.PlayerId})
		if awardError != nil {
			return errors.New("could not award track")
		}
	}

	createNextRoundError := gameHelpers.CreateNextRound(game.Id)
	if createNextRoundError != nil {
		return errors.New("could not create next round")
	}

	broadcastError := BroadcastGameUpdate(client, pool)
	if broadcastError != nil {
		return errors.New("could not broadcast game update")
	}

	return nil
}
