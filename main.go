package main

import (
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	a := App{}
	a.Initialize(
		os.Getenv("APP_DB_HOST"),
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	a.Run(":8080")
}