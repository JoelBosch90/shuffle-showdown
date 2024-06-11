package spotify

import (
	spotifyModels "api/lib/spotify/models"
	"encoding/json"
	"net/http"
	"net/url"
)

func RequestPlaylistInfo(playlistId string, countryCode string) (spotifyModels.Playlist, error) {
	path := "v1/playlists/" + url.QueryEscape(playlistId)

	headers := []Header{}
	params := []Param{
		{Name: "market", Value: countryCode},
		{Name: "fields", Value: GetSpotifyModelFields(spotifyModels.Playlist{})},
	}

	playlistResponse, playlistRequestError := ApiRequest(http.MethodGet, path, headers, params)
	if playlistRequestError != nil {
		return spotifyModels.Playlist{}, playlistRequestError
	}

	// Parse the response
	var playlistInfo spotifyModels.Playlist
	playlistDecoder := json.NewDecoder(playlistResponse.Body)
	decodeError := playlistDecoder.Decode(&playlistInfo)
	if decodeError != nil {
		return spotifyModels.Playlist{}, decodeError
	}

	if playlistInfo.Tracks.Limit >= playlistInfo.Tracks.Total {
		return playlistInfo, nil
	}

	// Get the next page of tracks
	additionalTrackItems, additionalTracksError := AddAdditionalTracks(&playlistInfo, path, headers, params)
	if additionalTracksError != nil {
		return spotifyModels.Playlist{}, additionalTracksError
	}
	playlistInfo.Tracks.Items = append(playlistInfo.Tracks.Items, additionalTrackItems...)

	return playlistInfo, nil
}
