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
	SentAt           time.Time `json:"sentAt"`
	FinishedLoading  bool      `json:"finishedLoading"`
	FinishedChecking bool      `json:"finishedChecking"`
	TracksChecked    uint      `json:"tracksChecked"`
	TracksVerified   uint      `json:"tracksVerified"`
	TracksLoaded     uint      `json:"tracksLoaded"`
	TracksTotal      uint      `json:"tracksTotal"`
}

type PlaylistCheckingCount struct {
	TracksChecked  uint `json:"tracks_checked"`
	TracksVerified uint `json:"tracks_verified"`
	TracksLoaded   uint `json:"tracks_loaded"`
}

type PlaylistTotal struct {
	TracksTotal uint `json:"tracks_total"`
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

	sendUpdate := func() {
		counts := PlaylistCheckingCount{}
		countsError := database.Raw(`
			SELECT
			  SUM(CASE WHEN check_completed_at IS NOT NULL THEN 1 ELSE 0 END) AS tracks_checked,
				SUM(CASE WHEN verified_at IS NOT NULL THEN 1 ELSE 0 END) AS tracks_verified,
				COUNT(*) AS tracks_loaded
			FROM tracks
			JOIN playlist_tracks ON tracks.id = playlist_tracks.track_id
			WHERE playlist_tracks.playlist_id = ?
		`, game.PlaylistId).Scan(&counts).Error
		if countsError != nil {
			return
		}

		tracksTotal := PlaylistTotal{}
		totalError := database.Raw(`SELECT tracks_total FROM playlists WHERE id = ?`, game.PlaylistId).Scan(&tracksTotal).Error
		if totalError != nil {
			return
		}

		pool.Broadcast <- ServerMessage{
			Type: ServerMessageTypePlaylistLoadingUpdate,
			Payload: PlaylistLoadingUpdate{
				SentAt:           time.Now(),
				FinishedLoading:  counts.TracksChecked > 0,
				FinishedChecking: counts.TracksChecked == counts.TracksLoaded,
				TracksChecked:    uint(counts.TracksChecked),
				TracksVerified:   uint(counts.TracksVerified),
				TracksLoaded:     uint(counts.TracksLoaded),
				TracksTotal:      uint(tracksTotal.TracksTotal),
			},
			GameId: client.GameId,
		}
	}

	playlistError := spotify.LoadFreshPlaylist(game.PlaylistId, game.CountryCode, sendUpdate)
	if playlistError != nil {
		client.SendError("Error loading playlist")
		return
	}

	go verification.VerifyTracks(sendUpdate)

	if !gameHelpers.HasEnoughTracks(client.GameId) {
		client.SendError("Too few tracks in playlist")
		return
	}
}
