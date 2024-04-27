package spotify

import (
	"api/database"
	"api/database/models"
	"api/lib/security"
	"time"

	uuid "github.com/satori/go.uuid"
)

const BUFFER_SECONDS int64 = 60

func getTokenFromDatabase() (models.AccessToken, error) {
	var token models.AccessToken
	database := database.Get()

	databaseError := database.Order("expires_at DESC").First(&token).Error
	if databaseError != nil || token.AccessToken == "" {
		return models.AccessToken{}, databaseError
	}

	decryptedToken, decryptionError := security.Decrypt(token.AccessToken)
	if decryptionError != nil {
		return models.AccessToken{}, decryptionError
	}

	bufferedNow := time.Unix(time.Now().Unix()+BUFFER_SECONDS, 0)
	if token.ExpiresAt.Before(bufferedNow) {
		return models.AccessToken{}, nil
	}

	return models.AccessToken{
		AccessToken: decryptedToken,
		TokenType:   token.TokenType,
		ExpiresAt:   token.ExpiresAt,
	}, nil
}

func storeTokenInDatabase(token models.AccessToken) error {
	database := database.Get()

	encryptedToken, encryptionError := security.Encrypt(token.AccessToken)
	if encryptionError != nil {
		return encryptionError
	}

	encryptedAccessToken := models.AccessToken{
		Id:          uuid.NewV4(),
		AccessToken: encryptedToken,
		TokenType:   token.TokenType,
		ExpiresAt:   token.ExpiresAt,
	}

	return database.Create(&encryptedAccessToken).Error
}

func GetAccessToken() (models.AccessToken, error) {
	// Ignore database errors as there might not be any tokens in the database.
	databaseToken, _ := getTokenFromDatabase()
	if databaseToken.AccessToken != "" {
		return databaseToken, nil
	}

	requestedToken, requestedTokenError := RequestNewAccessToken()
	if requestedTokenError != nil {
		return models.AccessToken{}, requestedTokenError
	}

	storeError := storeTokenInDatabase(requestedToken)
	return requestedToken, storeError
}
