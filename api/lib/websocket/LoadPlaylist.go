package websocket

import (
	"api/database"
	"api/database/models"
	gameHelpers "api/lib/game"
	"api/lib/spotify"
	"api/lib/verification"
	"time"
)

type PlaylistLoadingUpdate struct {
	SentAt       time.Time `json:"sentAt"`
	TracksLoaded uint      `json:"tracksLoaded"`
	TracksTotal  uint      `json:"tracksTotal"`
}

type PlaylistCheckingUpdate struct {
	SentAt         time.Time `json:"sentAt"`
	TracksChecked  uint      `json:"tracksChecked"`
	TracksVerified uint      `json:"tracksVerified"`
	TracksTotal    uint      `json:"tracksTotal"`
}

type PlaylistCheckingCount struct {
	TracksChecked  uint `json:"tracks_checked"`
	TracksVerified uint `json:"tracks_verified"`
	TracksTotal    uint `json:"tracks_total"`
}

func LoadPlaylist(client *Client) {
	database := database.Get()
	gameId := client.GameId
	game := models.Game{}
	pool := client.Pool

	databaseError := database.Where("id = ?", gameId).First(&game).Error
	if databaseError != nil {
		client.SendError("Game not found")
		return
	}

	sendLoadingUpdate := func(tracksLoaded int, totalTracks int) {
		pool.Broadcast <- ServerMessage{
			Type: ServerMessageTypePlaylistLoadingUpdate,
			Payload: PlaylistLoadingUpdate{
				SentAt:       time.Now(),
				TracksLoaded: uint(tracksLoaded),
				TracksTotal:  uint(totalTracks),
			},
			GameId: client.GameId,
		}
	}

	sendCheckingUpdate := func() {
		counts := PlaylistCheckingCount{}
		countsError := database.Raw(`
			SELECT
			  SUM(CASE WHEN check_completed_at IS NOT NULL THEN 1 ELSE 0 END) AS tracks_checked,
				SUM(CASE WHEN verified_at IS NOT NULL THEN 1 ELSE 0 END) AS tracks_verified,
				COUNT(*) AS tracks_total
			FROM tracks
			JOIN playlist_tracks ON tracks.id = playlist_tracks.track_id
			WHERE playlist_tracks.playlist_id = ?
		`, game.PlaylistId).Scan(&counts).Error
		if countsError != nil {
			return
		}

		pool.Broadcast <- ServerMessage{
			Type: ServerMessageTypePlaylistLoadingUpdate,
			Payload: PlaylistCheckingUpdate{
				SentAt:         time.Now(),
				TracksChecked:  uint(counts.TracksChecked),
				TracksVerified: uint(counts.TracksVerified),
				TracksTotal:    uint(counts.TracksTotal),
			},
			GameId: client.GameId,
		}
	}

	playlistError := spotify.LoadFreshPlaylist(game.PlaylistId, game.CountryCode, sendLoadingUpdate)
	if playlistError != nil {
		client.SendError("Error loading playlist")
		return
	}

	go verification.VerifyTracks(sendCheckingUpdate)

	if !gameHelpers.HasEnoughTracks(client.GameId) {
		client.SendError("Too few tracks in playlist")
		return
	}

	pool.Broadcast <- ServerMessage{
		Type: ServerMessageTypePlaylistLoaded,
		Payload: ErrorMessagePayload{
			Message: "Playlist loaded",
		},
		GameId: client.GameId,
	}
}
