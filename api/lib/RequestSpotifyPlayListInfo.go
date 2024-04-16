package lib

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

func RequestSpotifyPlayListInfo(playListId string) (SpotifyPlayListInfo, error) {
	path := "v1/playlists/" + url.QueryEscape(playListId)

	headers := []Header{}
	params := []Param{{Name: "market", Value: "NL"}}
	response, requestError := SpotifyAPIRequest(http.MethodGet, path, headers, params)
	if requestError != nil {
		return SpotifyPlayListInfo{}, requestError
	}

	// Parse the response
	var info SpotifyPlayListInfo
	decoder := json.NewDecoder(response.Body)
	decodeError := decoder.Decode(&info)
	if decodeError != nil {
		return SpotifyPlayListInfo{}, decodeError
	}

	log.Println(info)

	return info, nil
}
