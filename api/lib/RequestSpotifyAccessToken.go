package lib

import(
	"net/http"
	"encoding/json"
	"os"
	"time"
)

type SpotifyAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int64 `json:"expires_in"`
}

type ExtendedSpotifyAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresAt int64 `json:"expires_at"`
}

const SPOTIFY_TOKEN_REQUEST_URL = "https://accounts.spotify.com/api/token"

var cachedToken ExtendedSpotifyAccessToken

func RequestSpotifyAccessToken() (ExtendedSpotifyAccessToken, error) {
	if cachedToken.AccessToken != "" && cachedToken.ExpiresAt < time.Now().Unix() {
		return cachedToken, nil
	}

	// Create a new HTTP request
	request, requestError := http.NewRequest(http.MethodPost, SPOTIFY_TOKEN_REQUEST_URL, nil)
	if requestError != nil {
		return ExtendedSpotifyAccessToken{}, requestError
	}

	// Add the required headers
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Add the query parameters
	query := request.URL.Query()
	query.Add("grant_type", "client_credentials")
	query.Add("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
	query.Add("client_secret", os.Getenv("SPOTIFY_CLIENT_SECRET"))
	request.URL.RawQuery = query.Encode()

	// Send the request
	client := &http.Client{}
	response, responseError := client.Do(request)
	if responseError != nil {
		return ExtendedSpotifyAccessToken{}, responseError
	}

	// Parse the response
	var accessToken SpotifyAccessToken
	decoder := json.NewDecoder(response.Body)
	decodeError := decoder.Decode(&accessToken)
	if decodeError != nil {
		return ExtendedSpotifyAccessToken{}, decodeError
	}

	// Calculate the expiration time
	expirationUnixTimestamp := time.Now().Unix() + accessToken.ExpiresIn

	// Return the access token
	return ExtendedSpotifyAccessToken{
		AccessToken: accessToken.AccessToken,
		ExpiresAt: expirationUnixTimestamp,
	}, nil
}