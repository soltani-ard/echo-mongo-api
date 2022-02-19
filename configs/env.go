package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// EnvMongoURI => return mongodb uri with .env file
func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	// return uri
	return os.Getenv("MONGOURI")
}
