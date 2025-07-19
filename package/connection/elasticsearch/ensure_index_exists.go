package elasticsearch

import (
	"context"
	"fmt"
	"news-service/package/utils"
	"os"
	"strings"
)

func (es *Elasticsearch) EnsureIndexExistsFromFile(ctx context.Context, indexName string, fileNameDoc string) error {
	res, err := es.esClient.Indices.Exists([]string{indexName})
	if err != nil {
		return fmt.Errorf("failed to check index existence: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		fmt.Println("[Elasticsearch] Index", indexName, "already exists.")
		return nil
	}

	filePathName := utils.GetMappingPath(fileNameDoc)
	fmt.Println("log index file", filePathName)
	if res.StatusCode == 404 {
		// Baca file mapping
		mappingBytes, err := os.ReadFile(filePathName)
		if err != nil {
			return fmt.Errorf("failed to read mapping file: %w", err)
		}
		// Buat index dengan mapping dari file
		createRes, err := es.esClient.Indices.Create(
			indexName,
			es.esClient.Indices.Create.WithBody(strings.NewReader(string(mappingBytes))),
		)
		if err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
		defer createRes.Body.Close()

		if createRes.IsError() {
			return fmt.Errorf("create index error: %s", createRes.String())
		}

		fmt.Println("[Elasticsearch] Index", indexName, "created with mapping from file.")
		return nil
	}
	return fmt.Errorf("unexpected status code: %d", res.StatusCode)
}
