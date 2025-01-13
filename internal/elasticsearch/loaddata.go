package elasticsearch

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"go_day_03/internal/models"
	"os"
	"strconv"

	"github.com/elastic/go-elasticsearch/esapi"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

func LoadRestaurants(es *elasticsearch.Client, filePath string) error {
	exists, err := indexExists(es, "restaurants")
	if err != nil {
		return fmt.Errorf("error checking if index exists: %w", err)
	}

	if !exists {
		createIndexReq := esapi.IndicesCreateRequest{
			Index: "restaurants",
			Body: bytes.NewReader([]byte(`{
				"mappings": {
					"properties": {
						"Name": { "type": "text" },
						"Address": { "type": "text" },
						"Phone": { "type": "keyword" },
						"Longitude": { "type": "float" },
						"Latitude": { "type": "float" }
					}
				}
			}`)),
		}

		// next version //
		/*
			createIndexReq := esapi.IndicesCreateRequest{
				Index: "restaurants",
				Body: bytes.NewReader([]byte(`{
					"mappings": {
						"properties": {
							"Name": { "type": "text" },
							"Address": { "type": "text" },
							"Phone": { "type": "keyword" },
							"Location": { "type": geo_point}
						}
					}
				}`)),
			}
		*/

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
			ID:        record[0],
			Name:      record[1],
			Address:   record[2],
			Phone:     record[3],
			Longitude: longitude,
			Latitude:  latitude,
		}

		// for load most bigest data //
		meta := fmt.Sprintf(`{ "index": { "_id": "%s" } }%s`, restaurant.ID, "\n")
		buf.WriteString(meta)

		restaurantJSON, _ := json.Marshal(restaurant)
		buf.Write(restaurantJSON)
		buf.WriteString("\n")
	}

	bulkReq := esapi.BulkRequest{
		Index:   "restaurants",
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
