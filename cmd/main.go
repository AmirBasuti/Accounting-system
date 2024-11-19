package main

import (
	"AccountingSystem/internal/database"
)

func main() {
	if err := database.Connect(); err != nil {
		panic(err)
	}
	//database.Migrate()
}
