package spotify

import (
	spotifyModels "api/lib/spotify_models"
	"encoding/json"
	"net/http"
	"net/url"
)

func RequestPlayListInfo(playListId string, countryCode string) (spotifyModels.PlayList, error) {
	path := "v1/playlists/" + url.QueryEscape(playListId)

	headers := []Header{}
	params := []Param{{Name: "market", Value: countryCode}}
	playListResponse, playListRequestError := ApiRequest(http.MethodGet, path, headers, params)
	if playListRequestError != nil {
		return spotifyModels.PlayList{}, playListRequestError
	}

	// Parse the response
	var playListInfo spotifyModels.PlayList
	playListDecoder := json.NewDecoder(playListResponse.Body)
	decodeError := playListDecoder.Decode(&playListInfo)
	if decodeError != nil {
		return spotifyModels.PlayList{}, decodeError
	}

	if playListInfo.Tracks.Limit >= playListInfo.Tracks.Total {
		return playListInfo, nil
	}

	// Get the next page of tracks
	additionalTrackItems, additionalTracksError := AddAdditionalTracks(&playListInfo, path, headers, params)
	if additionalTracksError != nil {
		return spotifyModels.PlayList{}, additionalTracksError
	}
	playListInfo.Tracks.Items = append(playListInfo.Tracks.Items, additionalTrackItems...)

	return playListInfo, nil
}
