package spotify

import (
	"api/database/models"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

type SpotifyAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

const SPOTIFY_TOKEN_REQUEST_URL = "https://accounts.spotify.com/api/token"

func RequestNewAccessToken() (models.AccessToken, error) {
	// Create a new HTTP request
	request, requestError := http.NewRequest(http.MethodPost, SPOTIFY_TOKEN_REQUEST_URL, nil)
	if requestError != nil {
		return models.AccessToken{}, requestError
	}

	// Get the client ID and client secret from the environment
	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	if clientId == "" {
		return models.AccessToken{}, errors.New("missing Spotify client ID")
	}
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if clientSecret == "" {
		return models.AccessToken{}, errors.New("missing Spotify client secret")
	}

	// Add the required headers
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Add the query parameters
	query := request.URL.Query()
	query.Add("grant_type", "client_credentials")
	query.Add("client_id", clientId)
	query.Add("client_secret", clientSecret)
	request.URL.RawQuery = query.Encode()

	// Send the request
	client := &http.Client{}
	response, responseError := client.Do(request)
	if responseError != nil {
		return models.AccessToken{}, responseError
	}

	// Parse the response
	var accessToken SpotifyAccessToken
	decoder := json.NewDecoder(response.Body)
	decodeError := decoder.Decode(&accessToken)
	if decodeError != nil {
		return models.AccessToken{}, decodeError
	}

	// Calculate the expiration time
	expirationUnixTimestamp := time.Now().Unix() + accessToken.ExpiresIn

	// Return the access token
	return models.AccessToken{
		AccessToken: accessToken.AccessToken,
		TokenType:   accessToken.TokenType,
		ExpiresAt:   time.Unix(expirationUnixTimestamp, 0),
	}, nil
}
