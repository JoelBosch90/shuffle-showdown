package game

import "api/database/models"

type PlayerReturnObject struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func GetPlayerReturnObject(player models.Player) PlayerReturnObject {
	return PlayerReturnObject{
		Id:   player.Id.String(),
		Name: player.Name,
	}
}
