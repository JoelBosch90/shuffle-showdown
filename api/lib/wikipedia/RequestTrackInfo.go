package wikipedia

import (
	languageModels "api/lib/wikipedia/languages/models"
	wikipediaModels "api/lib/wikipedia/models"
	"encoding/json"
	"net/http"
)

func RequestTrackInfo(language languageModels.Language, trackTitle string) (wikipediaModels.Response, error) {
	var formattedResponse wikipediaModels.Response

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
