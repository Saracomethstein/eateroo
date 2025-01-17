package handlers

import (
	"go_day_03/internal/elasticsearch"
	"log"
	"net/http"
	"strconv"

	esearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/labstack/echo/v4"
)

func Ping(c echo.Context) error {
	return c.String(http.StatusOK, "ok\n")
}

func GetPlace(c echo.Context) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 || limit > 100 {
		limit = 30
	}

	es, err := esearch.NewClient(esearch.Config{
		Addresses: []string{"http://elasticsearch:9200"},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create Elasticsearch client"})
	}

	// Получение данных из Elasticsearch
	index := "places"
	restaurants, total, err := elasticsearch.FetchRestaurants(es, index, page, limit)
	if err != nil {
		log.Printf("error fetching restaurants: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch restaurants"})
	}

	response := map[string]interface{}{
		"page":        page,
		"limit":       limit,
		"total":       total,
		"restaurants": restaurants,
	}

	return c.JSON(http.StatusOK, response)
}
