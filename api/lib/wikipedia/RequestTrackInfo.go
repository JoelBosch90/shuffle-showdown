package wikipedia

import (
	languages "api/lib/wikipedia/languages"
	"encoding/json"
	"net/http"
)

func RequestTrackInfo(language languages.Language, trackTitle string) (Response, error) {
	var formattedResponse Response

	headers := []Header{}
	params := []Param{
		{Name: "action", Value: "query"},
		{Name: "prop", Value: "categories|links"},
		{Name: "titles", Value: trackTitle},
		{Name: "utf8", Value: "1"},
		{Name: "format", Value: "json"},
		{Name: "formatversion", Value: "2"},
		{Name: "cllimit", Value: "max"},
		{Name: "pllimit", Value: "max"},
		{Name: "redirects", Value: "1"},
	}

	response, requestError := ApiRequest(http.MethodGet, headers, params, language)
	if requestError != nil {
		return formattedResponse, requestError
	}

	decoder := json.NewDecoder(response.Body)
	decodeError := decoder.Decode(&formattedResponse)
	if decodeError != nil {
		return formattedResponse, decodeError
	}

	return formattedResponse, nil
}
