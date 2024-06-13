package database

import (
	"api/database/purge"

	"github.com/jasonlvhit/gocron"
)

var MAX_AGE_IN_DAYS int = 30

func purgeAll() error {
	database := Get()

	purgeAllError := purge.PurgeAll(database, MAX_AGE_IN_DAYS)
	if purgeAllError != nil {
		return purgeAllError
	}

	return nil
}

func PurgeDaily() {
	gocron.Every(1).Minute().Do(purgeAll)
	<-gocron.Start()
}
