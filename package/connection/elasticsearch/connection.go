package elasticsearch

import (
	"context"
	"fmt"
	"news-service/package/config"

	"github.com/elastic/go-elasticsearch/v7"
)

type Elasticsearch struct {
	esClient *elasticsearch.Client
	conf     *config.Config
}

type ElasticsearchClient interface {
	SearchInElasticsearch(ctx context.Context, index string, query map[string]interface{}, dest interface{}) error
	EnsureIndexExistsFromFile(ctx context.Context, indexName string, fileNameDoc string) error 
	StoreToElasticsearch(ctx context.Context, payload interface{}) error
}

func NewElasticSearch(conf *config.Config) (*Elasticsearch, error) {
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			conf.ElasticsearchURL,
		},
		Username: conf.ElasticsearchUsername,
		Password: conf.ElasticsearchPassword,
	})
	if err != nil {
		return &Elasticsearch{}, fmt.Errorf("cannot create Elasticsearch client: %v", err)
	}
	return &Elasticsearch{
		esClient: esClient,
		conf:     conf,
	}, nil
}