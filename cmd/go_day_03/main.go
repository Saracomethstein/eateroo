package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// scheme.json for struct index //
	mapping :=
		`
	{
		"mappings": {
			"properties": {	
				"name": {
					"type": "text"
				},
				"address": {
					"type": "text"
				}, 
				"phone": {
					"type": "text"
				},
				"location": {
					"type": "geo_point"
				}
			}
		}
	}
	`

	req := esapi.IndicesCreateRequest{
		Index: "places",
		Body:  bytes.NewReader([]byte(mapping)),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error creating index: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error response: %s", res.String())
	}

	fmt.Println("Index created succesfully.")
}
