package service

import (
	"errors"
	"fmt"
	"go_day_03/internal/elasticsearch"
	"go_day_03/internal/models"
	"go_day_03/internal/repositories"

	esearch "github.com/elastic/go-elasticsearch/v8"
)

type Store interface {
	GetPlace(limit, offset int) ([]models.Restaurant, int, error)
}

func GetPlace(page, limit int, searchQuery string) ([]models.Restaurant, int, error) {
	config := repositories.New()
	es, err := esearch.NewClient(esearch.Config{
		Addresses: []string{config.ESearchURL},
	})
	if err != nil {
		return nil, 0, errors.New(fmt.Sprintf("failed to create Elasticsearch client"))
	}

	restaurants, total, err := elasticsearch.FetchRestaurants(es, page, limit, searchQuery)
	if err != nil {
		return nil, 0, errors.New(fmt.Sprintf("failed to fetch restaurants"))
	}

	return restaurants, total, nil
}
