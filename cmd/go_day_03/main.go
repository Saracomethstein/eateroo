package main

import (
	"fmt"
	"go_day_03/internal/elasticsearch"
	"log"
	"os"

	esearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/joho/godotenv"
)

func main() {
	es, err := esearch.NewClient(esearch.Config{
		Addresses: []string{"http://elasticsearch:9200"},
	})
	if err != nil {
		fmt.Printf("Error connecting to Elasticsearch: %s\n", err)
		return
	}

	if err := elasticsearch.LoadRestaurants(es, getEnv()); err != nil {
		fmt.Printf("Error loading restaurants: %s\n", err)
	}
}

func getEnv() string {
	if err := godotenv.Load("/app/.env"); err != nil {
		log.Println("Error: ", err)
	}

	sourceSCV := os.Getenv("DATA_SOURCE")

	if sourceSCV == "" {
		log.Fatal("DATA_SOURCE environment variable are missing.")
	}

	return sourceSCV
}
