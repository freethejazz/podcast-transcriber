package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type Caption struct {
	Url           string        `json:"url"`
	Index         int           `json:"index"`
	Text          string        `json:"text"`
	Context       string        `json:"context"`
	TimestampFrom time.Duration `json:"timestamp_from"`
	TimestampTo   time.Duration `json:"timestamp_to"`
	ClipLength    time.Duration `json:"clip_length"`
}

func SearchCaptions(search string) ([]Caption, error) {
	var r map[string]interface{}

	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"}, // Update with your Elasticsearch host
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %v", err)
	}

	indexName := "captions"

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"context": search,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex(indexName),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	hits, ok := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid search result format")
	}

	captions := make([]Caption, 0)

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]

		captionBytes, err := json.Marshal(source)
		if err != nil {
			log.Printf("Failed to marshal caption source: %s", err)
			continue
		}

		var caption Caption
		if err := json.Unmarshal(captionBytes, &caption); err != nil {
			log.Printf("Failed to unmarshal caption: %s", err)
			continue
		}

		captions = append(captions, caption)
	}

	return captions, nil
}

func IndexCaptions(captions []Caption) error {
	// Create Elasticsearch client
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"}, // Update with your Elasticsearch host
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create Elasticsearch client: %v", err)
	}

	indexName := "captions"

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  indexName,
		Client: client,
	})

	// Prepare the bulk update request
	for _, caption := range captions {
		// Serialize the caption struct
		captionBytes, err := json.Marshal(caption)
		if err != nil {
			return fmt.Errorf("failed to serialize caption struct: %v", err)
		}

		addErr := bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action: "index",
				Body:   bytes.NewReader(captionBytes),
			},
		)
		if addErr != nil {
			log.Fatalf("Unexpected indexing err: %s", err)
		}

	}

	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	log.Println("Captions indexed successfully.")
	return nil
}
