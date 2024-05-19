package game

import "api/database/models"

func IncludesPlayer(haystack []models.Player, needle models.Player) bool {
	for _, player := range haystack {
		if player.Id == needle.Id {
			return true
		}
	}
	return false
}
