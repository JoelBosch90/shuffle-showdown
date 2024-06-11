package spotify

import (
	"log"
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

const SPOTIFY_API_URL = "https://api.spotify.com/"

func ApiRequest(method string, path string, headers []Header, params []Param) (*http.Response, error) {
	url := SPOTIFY_API_URL + path

	log.Println("URL: ", url)
	log.Println("PARAMS: ", params)

	// Create a new HTTP request
	request, requestError := http.NewRequest(method, url, nil)
	if requestError != nil {
		return nil, requestError
	}

	// Get an access token
	accessToken, accessTokenError := GetAccessToken()
	if accessTokenError != nil {
		return nil, accessTokenError
	}

	// Add the authorization header.
	request.Header.Add("Authorization", accessToken.TokenType+" "+accessToken.AccessToken)

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
