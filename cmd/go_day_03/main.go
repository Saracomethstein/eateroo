package main

import (
	"fmt"
	"go_day_03/internal/elasticsearch"
	"log"
	"net/http"
	"os"

	esearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/joho/godotenv"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ping server //
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok\n")
	})

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

	e.Logger.Fatal(e.Start(":8888"))
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
