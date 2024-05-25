package websocket

import (
	"api/database"
	"api/database/models"
	gameHelpers "api/lib/game"
	"encoding/json"
	"errors"
	"log"
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

	log.Println("CURRENT ROUND: ", currentRound.Track)

	var answer Answer
	answerParseError := json.Unmarshal([]byte(message.Payload), &answer)
	if answerParseError != nil {
		return errors.New("could not parse answer")
	}
	if answer.BeforeReleaseYear == nil && answer.AfterReleaseYear == nil {
		return errors.New("answer must contain at least one field")
	}

	if answer.BeforeReleaseYear != nil {
		log.Println("BEFORE RELEASE YEAR: ", *answer.BeforeReleaseYear)
	}
	if answer.AfterReleaseYear != nil {
		log.Println("AFTER RELEASE YEAR: ", *answer.AfterReleaseYear)
	}

	correctBefore := answer.BeforeReleaseYear == nil || currentRound.Track.ReleaseYear <= uint(*answer.BeforeReleaseYear)
	correctAfter := answer.AfterReleaseYear == nil || currentRound.Track.ReleaseYear >= uint(*answer.AfterReleaseYear)
	if correctBefore && correctAfter {
		awardError := gameHelpers.AwardTrack(game.Id, currentRound.Track, models.Player{Id: client.PlayerId})
		if awardError != nil {
			return errors.New("could not award track")
		}
	}
	log.Println("CORRECT BEFORE: ", correctBefore)
	log.Println("CORRECT AFTER: ", correctAfter)

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
