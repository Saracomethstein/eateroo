package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"go_day_03/internal/models"
	"os"
	"strconv"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func main() {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://elasticsearch:9200"},
	})
	if err != nil {
		fmt.Printf("Ошибка подключения к Elasticsearch: %s\n", err)
		return
	}

	file, err := os.Open("data/data.csv")
	if err != nil {
		fmt.Printf("Ошибка при открытии файла: %s\n", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.LazyQuotes = true
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Ошибка при чтении CSV: %s\n", err)
		return
	}

	if len(records) < 2 {
		fmt.Println("CSV-файл пуст или не содержит данных")
		return
	}

	req := esapi.IndicesCreateRequest{
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
	res, err := req.Do(context.Background(), es)
	if err != nil || res.IsError() {
		fmt.Printf("Ошибка при создании индекса: %s\n", res.String())
		return
	}
	defer res.Body.Close()

	for i, record := range records[1:] {
		longitude, _ := strconv.ParseFloat(record[4], 64)
		latitude, _ := strconv.ParseFloat(record[5], 64)

		restaurant := models.Restaurant{
			ID:        record[0],
			Name:      record[1],
			Address:   record[2],
			Phone:     record[3],
			Longitude: longitude,
			Latitude:  latitude,
		}

		data, err := json.Marshal(restaurant)
		if err != nil {
			fmt.Printf("Ошибка при преобразовании в JSON на строке %d: %s\n", i+2, err)
			continue
		}

		req := esapi.IndexRequest{
			Index:      "restaurants",
			DocumentID: restaurant.ID,
			Body:       bytes.NewReader(data),
			Refresh:    "true",
		}

		res, err := req.Do(context.Background(), es)
		if err != nil || res.IsError() {
			fmt.Printf("Ошибка при загрузке данных на строке %d: %s\n", i+2, res.String())
			continue
		}
		defer res.Body.Close()

		fmt.Printf("Успешно загружен ресторан с ID %s Name: %s\n", restaurant.ID, restaurant.Name)
	}
}
