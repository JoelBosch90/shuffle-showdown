package lib

import (
	"api/database"
	"api/database/models"
	"time"
)

const BUFFER_SECONDS int64 = 60

func getTokenFromDatabase() models.AccessToken {
	var token models.AccessToken
	database := database.Get()

	databaseError := database.Order("expires_at DESC").First(&token).Error
	decryptedToken, decryptionError := Decrypt(token.AccessToken)

	minuteAgo := time.Unix(time.Now().Unix()-BUFFER_SECONDS, 0)

	if databaseError == nil && decryptionError == nil && token.ExpiresAt.After(minuteAgo) {
		return models.AccessToken{
			AccessToken: decryptedToken,
			TokenType:   token.TokenType,
			ExpiresAt:   token.ExpiresAt,
		}
	}

	return models.AccessToken{}
}

func storeTokenInDatabase(token models.AccessToken) {
	database := database.Get()

	encryptedToken, encryptionError := Encrypt(token.AccessToken)
	if encryptionError != nil {
		return
	}

	encryptedAccessToken := models.AccessToken{
		AccessToken: encryptedToken,
		TokenType:   token.TokenType,
		ExpiresAt:   token.ExpiresAt,
	}

	database.Create(&encryptedAccessToken)
}

func GetSpotifAccessToken() (models.AccessToken, error) {
	token := getTokenFromDatabase()
	if token.AccessToken != "" {
		return token, nil
	}

	token, tokenError := RequestNewSpotifyAccessToken()
	if tokenError != nil {
		return models.AccessToken{}, tokenError
	}

	storeTokenInDatabase(token)
	return token, nil
}
