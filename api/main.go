package main

import (
	"api/database"
	"api/router"
)

func main() {
	go database.PurgeDaily()
	router.Run()
}
