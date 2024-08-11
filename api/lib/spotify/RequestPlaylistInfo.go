package spotify

import (
	spotifyModels "api/lib/spotify/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

func RequestPlaylistInfo(playlistId string, countryCode string, sendUpdate func(isLiveUpdate bool)) (spotifyModels.Playlist, error) {
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
	if playlistInfo.Id == "" {
		return spotifyModels.Playlist{}, errors.New("playlist empty")
	}

	sendUpdate(true)
	if playlistInfo.Tracks.Limit >= playlistInfo.Tracks.Total {
		return playlistInfo, nil
	}

	sendAdditionalUpdate := func(additionalTracksLoaded int) {
		sendUpdate(true)
	}

	// Get the next page of tracks
	additionalTrackItems, additionalTracksError := AddAdditionalTracks(&playlistInfo, path, headers, params, sendAdditionalUpdate)
	if additionalTracksError != nil {
		return spotifyModels.Playlist{}, additionalTracksError
	}
	playlistInfo.Tracks.Items = append(playlistInfo.Tracks.Items, additionalTrackItems...)

	sendUpdate(true)
	return playlistInfo, nil
}
