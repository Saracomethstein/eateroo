package main

import (
	"fmt"
	"go_day_03/internal/elasticsearch"
	"go_day_03/internal/handlers"
	"go_day_03/internal/repositories"

	esearch "github.com/elastic/go-elasticsearch/v8"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	Load()

	e.GET("/api/ping", handlers.Ping)
	e.GET("/api/places", handlers.GetPlace)

	e.Logger.Fatal(e.Start(":8888"))
}

func Load() {
	config := repositories.New()
	es, err := esearch.NewClient(esearch.Config{
		Addresses: []string{config.ESearchURL},
	})

	if err != nil {
		fmt.Printf("Error connecting to Elasticsearch: %s\n", err)
		return
	}

	if err := elasticsearch.LoadRestaurants(es, config.DataCSV); err != nil {
		fmt.Printf("Error loading restaurants: %s\n", err)
	}
}
