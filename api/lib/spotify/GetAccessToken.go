package spotify

import (
	"api/database"
	"api/database/models"
	"api/lib/security"
	"time"

	uuid "github.com/satori/go.uuid"
)

const BUFFER_SECONDS int64 = 60

func getTokenFromDatabase() (models.SpotifyAccessToken, error) {
	var token models.SpotifyAccessToken
	database := database.Get()

	databaseError := database.Order("expires_at DESC").First(&token).Error
	if databaseError != nil || token.AccessToken == "" {
		return models.SpotifyAccessToken{}, databaseError
	}

	decryptedToken, decryptionError := security.Decrypt(token.AccessToken)
	if decryptionError != nil {
		return models.SpotifyAccessToken{}, decryptionError
	}

	bufferedNow := time.Unix(time.Now().Unix()+BUFFER_SECONDS, 0)
	if token.ExpiresAt.Before(bufferedNow) {
		return models.SpotifyAccessToken{}, nil
	}

	return models.SpotifyAccessToken{
		AccessToken: decryptedToken,
		TokenType:   token.TokenType,
		ExpiresAt:   token.ExpiresAt,
	}, nil
}

func storeTokenInDatabase(token models.SpotifyAccessToken) error {
	database := database.Get()

	encryptedToken, encryptionError := security.Encrypt(token.AccessToken)
	if encryptionError != nil {
		return encryptionError
	}

	encryptedAccessToken := models.SpotifyAccessToken{
		Id:          uuid.NewV4(),
		AccessToken: encryptedToken,
		TokenType:   token.TokenType,
		ExpiresAt:   token.ExpiresAt,
	}

	return database.Create(&encryptedAccessToken).Error
}

func GetAccessToken() (models.SpotifyAccessToken, error) {
	// Ignore database errors as there might not be any tokens in the database.
	databaseToken, _ := getTokenFromDatabase()
	if databaseToken.AccessToken != "" {
		return databaseToken, nil
	}

	requestedToken, requestedTokenError := RequestNewAccessToken()
	if requestedTokenError != nil {
		return models.SpotifyAccessToken{}, requestedTokenError
	}

	storeError := storeTokenInDatabase(requestedToken)
	return requestedToken, storeError
}
