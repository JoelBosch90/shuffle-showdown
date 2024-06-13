package main

import (
	"api/database"
	"api/router"
)

func main() {
	router.Run()
	go database.PurgeDaily()
}
