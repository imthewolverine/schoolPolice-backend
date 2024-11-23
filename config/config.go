package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
    // Attempt to load using both relative and absolute paths for robustness
    if err := godotenv.Load(".env"); err != nil {
        if err := godotenv.Load("/app/.env"); err != nil {
            log.Fatalf("Error loading .env file")
        }
    }
    fmt.Println("Loaded .env file")
}


func GetEnv(key string) string {
    return os.Getenv(key)
}
