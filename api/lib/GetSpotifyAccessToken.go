package lib

import (
	"api/database"
	"api/database/models"
	"time"
)

func getTokenFromDatabase() models.AccessToken {
	var token models.AccessToken
	database := database.Get()

	databaseError := database.Order("expires_at DESC").First(&token).Error
	decryptedToken, decryptionError := Decrypt(token.AccessToken)

	if databaseError == nil && decryptionError == nil && token.ExpiresAt.After(time.Now()) {
		return models.AccessToken{
			AccessToken: decryptedToken,
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
