package spotify

import (
	"api/database"
	"api/database/models"
	"api/lib/security"
	"time"
)

const BUFFER_SECONDS int64 = 60

func getTokenFromDatabase() models.AccessToken {
	var token models.AccessToken
	database := database.Get()

	databaseError := database.Order("expires_at DESC").First(&token).Error
	decryptedToken, decryptionError := security.Decrypt(token.AccessToken)

	bufferedNow := time.Unix(time.Now().Unix()+BUFFER_SECONDS, 0)
	if databaseError == nil && decryptionError == nil && token.ExpiresAt.After(bufferedNow) {
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

	encryptedToken, encryptionError := security.Encrypt(token.AccessToken)
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

func GetAccessToken() (models.AccessToken, error) {
	token := getTokenFromDatabase()
	if token.AccessToken != "" {
		return token, nil
	}

	token, tokenError := RequestNewAccessToken()
	if tokenError != nil {
		return models.AccessToken{}, tokenError
	}

	storeTokenInDatabase(token)
	return token, nil
}