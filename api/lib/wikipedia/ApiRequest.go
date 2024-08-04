package wikipedia

import (
	helpers "api/lib/helpers"
	languages "api/lib/wikipedia/languages"
	languageModels "api/lib/wikipedia/languages/models"
	"net/http"
)

type Header struct {
	Name  string
	Value string
}

type Param struct {
	Name  string
	Value string
}

func ApiRequest(method string, headers []Header, params []Param, language languageModels.Language) (*http.Response, error) {
	apiUrl := languages.Map[language].ApiUrl

	// Create a new HTTP request
	request, requestError := http.NewRequest(method, apiUrl, nil)
	if requestError != nil {
		return nil, requestError
	}

	request.Header.Add("User-Agent", helpers.UserAgent())

	// Add the other headers.
	for _, header := range headers {
		request.Header.Add(header.Name, header.Value)
	}

	// Add the query parameters
	query := request.URL.Query()

	for _, param := range params {
		query.Add(param.Name, param.Value)
	}
	request.URL.RawQuery = query.Encode()

	// Send the request
	client := &http.Client{}
	response, responseError := client.Do(request)
	if responseError != nil {
		return nil, responseError
	}

	return response, nil
}
