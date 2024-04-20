package spotify

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func RequestPlayListInfo(playListId string, countryCode string) (PlayListInfo, error) {
	path := "v1/playlists/" + url.QueryEscape(playListId)

	headers := []Header{}
	params := []Param{{Name: "market", Value: countryCode}}
	response, requestError := ApiRequest(http.MethodGet, path, headers, params)
	if requestError != nil {
		return PlayListInfo{}, requestError
	}

	// Parse the response
	var info PlayListInfo
	decoder := json.NewDecoder(response.Body)
	decodeError := decoder.Decode(&info)
	if decodeError != nil {
		return PlayListInfo{}, decodeError
	}

	return info, nil
}
