package main

import (
	"MagmaAPI/app"
	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	_ = godotenv.Load()

	app.Start()
}
