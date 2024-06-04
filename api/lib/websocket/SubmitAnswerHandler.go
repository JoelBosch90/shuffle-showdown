package websocket

import (
	"api/database"
	"api/database/models"
	gameHelpers "api/lib/game"
	"encoding/json"
	"errors"
)

type Answer struct {
	AfterReleaseYear  *int `json:"afterReleaseYear"`
	BeforeReleaseYear *int `json:"beforeReleaseYear"`
}

func SubmitAnswerHandler(message ClientMessage, client *Client, pool *ConnectionPool) error {
	database := database.Get()
	var game models.Game

	gameError := database.Preload("Rounds.Track").Preload("WonTracks.Track").Where("id = ?", client.GameId).First(&game).Error
	if gameError != nil {
		return errors.New("could not find game")
	}
	if game.HasFinished {
		return errors.New("game has already finished")
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

	correctAnswer, verifyError := gameHelpers.VerifyAnswer(answer.AfterReleaseYear, answer.BeforeReleaseYear, client.PlayerId, &game.Rounds, &game.WonTracks)
	if verifyError != nil {
		return verifyError
	}
	if correctAnswer {
		awardError := gameHelpers.AwardTrack(game.Id, currentRound.Track, models.Player{Id: client.PlayerId})
		if awardError != nil {
			return errors.New("could not award track")
		}
	}

	hasFinished := gameHelpers.CheckGameEnd(game.Id)
	if !hasFinished {
		createNextRoundError := gameHelpers.CreateNextRound(game.Id)
		if createNextRoundError != nil {
			return errors.New("could not create next round")
		}
	}

	broadcastError := BroadcastGameUpdate(client, pool)
	if broadcastError != nil {
		return errors.New("could not broadcast game update")
	}

	return nil
}
