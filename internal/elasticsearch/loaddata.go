package elasticsearch

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"go_day_03/internal/models"
	"log"
	"os"
	"strconv"

	"github.com/elastic/go-elasticsearch/esapi"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

func LoadRestaurants(es *elasticsearch.Client, filePath string) error {
	exists, err := indexExists(es, "places")
	if err != nil {
		return fmt.Errorf("error checking if index exists: %w", err)
	}

	if !exists {
		createIndexReq := esapi.IndicesCreateRequest{
			Index: "places",
			Body: bytes.NewReader([]byte(`{
					"mappings": {
						"properties": {
							"Name": { "type": "text" },
							"Address": { "type": "text" },
							"Phone": { "type": "keyword" },
							"Location": { "type": "geo_point"}
						}
					}
				}`)),
		}

		res, err := createIndexReq.Do(context.Background(), es)
		if err != nil {
			return fmt.Errorf("error creating index: %w", err)
		}
		defer res.Body.Close()
		if res.IsError() {
			return fmt.Errorf("error creating index: %s", res.String())
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.LazyQuotes = true
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file is empty or lacks data")
	}

	var buf bytes.Buffer
	for i, record := range records[1:] {
		longitude, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			fmt.Printf("Error parsing longitude on row %d: %s\n", i+2, err)
			continue
		}
		latitude, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			fmt.Printf("Error parsing latitude on row %d: %s\n", i+2, err)
			continue
		}

		restaurant := models.Restaurant{
			ID:      record[0],
			Name:    record[1],
			Address: record[2],
			Phone:   record[3],
			Location: struct {
				Longitude float64 `json: "Longitude"`
				Latitude  float64 `json: "Latitude"`
			}{
				Longitude: longitude,
				Latitude:  latitude,
			},
		}

		meta := fmt.Sprintf(`{ "index": { "_id": "%s" } }%s`, restaurant.ID, "\n")
		buf.WriteString(meta)
		restaurantJSON, _ := json.Marshal(restaurant)
		buf.Write(restaurantJSON)
		buf.WriteString("\n")

		fmt.Printf("Restaurant ID: %s, Name: %s, Phone: %s, Location:%f\t%f\n", restaurant.ID, restaurant.Name, restaurant.Phone, restaurant.Location.Longitude, restaurant.Location.Latitude)
	}

	bulkReq := esapi.BulkRequest{
		Index:   "places",
		Body:    &buf,
		Refresh: "true",
	}
	res, err := bulkReq.Do(context.Background(), es)
	if err != nil {
		return fmt.Errorf("error executing bulk request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk request error: %s", res.String())
	}

	fmt.Println("All records successfully uploaded.")
	return nil
}

func indexExists(es *elasticsearch.Client, indexName string) (bool, error) {
	req := esapi.IndicesExistsRequest{
		Index: []string{indexName},
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	return res.StatusCode == 200, nil
}

func FetchRestaurants(client *elasticsearch.Client, index string, page, limit int) ([]models.Restaurant, int, error) {
	from := (page - 1) * limit

	// Elasticsearch запрос
	query := map[string]interface{}{
		"from": from,
		"size": limit,
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, fmt.Errorf("error encoding query: %w", err)
	}

	// Выполнение запроса к Elasticsearch
	res, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex(index),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting response: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("error in response: %s", res.String())
	}

	var esResponse models.ElasticsearchResponse
	if err := json.NewDecoder(res.Body).Decode(&esResponse); err != nil {
		return nil, 0, fmt.Errorf("error parsing response body: %w", err)
	}

	restaurants := make([]models.Restaurant, 0, len(esResponse.Hits.Hits))
	for _, hit := range esResponse.Hits.Hits {
		var restaurant models.Restaurant
		if err := json.Unmarshal(hit.Source, &restaurant); err != nil {
			log.Printf("error unmarshaling hit source: %v", err)
			continue
		}
		restaurant.ID = hit.ID
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, esResponse.Hits.Total.Value, nil
}
