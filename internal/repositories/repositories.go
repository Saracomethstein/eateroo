package repositories

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DataCSV      string
	ESearchURL   string
	ESearchIndex string
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("Error: No .env file found.")
	}
}

func New() *Config {
	return &Config{
		DataCSV:      getEnv("DATA_SOURCE", ""),
		ESearchURL:   getEnv("ELASTICSEARCH_SOURCE", ""),
		ESearchIndex: getEnv("ELASTICSEARCH_INDEX", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
