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
	database.AutoMigrate(&models.Game{}, &models.PlayList{})

	return database
}

func addModels(database *gorm.DB) *gorm.DB {
	database.AutoMigrate(&models.Game{}, &models.PlayList{})

	return database
}

func Get() *gorm.DB {
	name := "shuffle"
	database := connect(name)

	addModels(database)

	return database
}