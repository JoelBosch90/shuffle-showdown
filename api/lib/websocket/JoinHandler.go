package websocket

import (
	gameHelpers "api/lib/game"
	"errors"
)

func JoinHandler(message ClientMessage, client *Client, pool *ConnectionPool) error {
	addError := gameHelpers.AddPlayerToGame(client.GameId, client.PlayerId)
	if addError != nil {
		return addError
	}

	awardInitialTracksError := gameHelpers.AwardInitialTrack(client.GameId, client.PlayerId)
	if awardInitialTracksError != nil {
		return errors.New("could not award initial track")
	}

	broadcastError := BroadcastGameUpdate(client, pool)
	if broadcastError != nil {
		return errors.New("could not broadcast game update")
	}

	return nil
}
