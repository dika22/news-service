package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

func (es Elasticsearch) StoreToElasticsearch(ctx context.Context, payload interface{}) error {

	if err := es.EnsureIndexExistsFromFile(ctx, es.conf.ArticleIndex, "articles.json"); err != nil {
		fmt.Println("Failed to ensure index exists:", err)
		return err
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	res, err := es.esClient.Index(
		es.conf.ArticleIndex,
		bytes.NewReader(requestBody),
		es.esClient.Index.WithContext(ctx),
		es.esClient.Index.WithDocumentType("_doc"), // opsional, tergantung versi ES
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("elasticsearch error: %s", body)
	}

	return nil
}