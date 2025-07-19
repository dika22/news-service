package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Handler for searching data from Elasticsearch
func (es Elasticsearch) SearchInElasticsearch(ctx context.Context, indexName string, query map[string]interface{}, dest interface{}) error {
	// Marshal query map ke JSON
	queryBody, err := json.Marshal(query)
	if err != nil {
		return err
	}

	// Kirim request ke Elasticsearch
	res, err := es.esClient.Search(
		es.esClient.Search.WithContext(ctx),
		es.esClient.Search.WithIndex(indexName),
		es.esClient.Search.WithBody(strings.NewReader(string(queryBody))),
		es.esClient.Search.WithTrackTotalHits(true),
		es.esClient.Search.WithPretty(),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("search error: %s", body)
	}


	// Decode hasil ke dest
	if err := json.NewDecoder(res.Body).Decode(&dest); err != nil {
		return err
	}
	return nil
}
