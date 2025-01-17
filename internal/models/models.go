package models

import "encoding/json"

type Restaurant struct {
	ID       string `json: "ID"`
	Name     string `json: "Name"`
	Address  string `json: "Address"`
	Phone    string `json: "Phone"`
	Location struct {
		Longitude float64 `json: "Longitude"`
		Latitude  float64 `json: "Latitude"`
	} `json: "Location"`
}

type ElasticsearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			ID     string          `json:"_id"`
			Source json.RawMessage `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
