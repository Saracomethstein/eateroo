package main

import (
	"fmt"
	"go_day_03/internal/elasticsearch"

	esearch "github.com/elastic/go-elasticsearch/v8"
)

func main() {
	es, err := esearch.NewClient(esearch.Config{
		Addresses: []string{"http://elasticsearch:9200"},
	})
	if err != nil {
		fmt.Printf("Error connecting to Elasticsearch: %s\n", err)
		return
	}

	// move to env file //
	filePath := "data/data.csv"
	if err := elasticsearch.LoadRestaurants(es, filePath); err != nil {
		fmt.Printf("Error loading restaurants: %s\n", err)
	}
}
