package game

import "api/database/models"

func FindNextRoundNumber(rounds []models.Round) uint {
	if len(rounds) == 0 {
		return 1
	}

	return rounds[len(rounds)-1].Number + 1
}
