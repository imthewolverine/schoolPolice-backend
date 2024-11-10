package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
    err := godotenv.Load(".env")
    
    fmt.Printf("Loaded .env file\n");
    if err != nil {
        log.Fatalf("Error loading .env file")
    }
}

func GetEnv(key string) string {
    return os.Getenv(key)
}
