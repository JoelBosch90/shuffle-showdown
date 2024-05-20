package game

import "api/database/models"

func LastRound(rounds []models.Round) (int, models.Round) {
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
