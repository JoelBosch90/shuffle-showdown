package database

import (
	"api/database/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func connect(name string) *gorm.DB {
	database, error := gorm.Open("sqlite3", name+".db")

	if error != nil {
		panic(error)
	}

	database.LogMode(true)

	return database
}

func addModels(database *gorm.DB) *gorm.DB {
	database.AutoMigrate(&models.Game{}, &models.Playlist{}, &models.AccessToken{}, &models.Artist{}, &models.Track{}, &models.Player{})

	return database
}

func Get() *gorm.DB {
	name := "shuffle"
	database := addModels(connect(name))

	return database
}
