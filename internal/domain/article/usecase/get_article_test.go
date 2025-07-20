package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/dika22/news-service/internal/domain/article/repository/mocks"
	authorMocks "github.com/dika22/news-service/internal/domain/author/repository/mocks"
	"github.com/dika22/news-service/package/config"
	mqMocks "github.com/dika22/news-service/package/rabbit-mq/mocks"
	"github.com/dika22/news-service/package/structs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cacheMock "github.com/dika22/news-service/internal/domain/article/repository/cache/mocks"
	esMocks "github.com/dika22/news-service/package/connection/elasticsearch/mocks"
)

// -------------------- TEST ------------------------
func TestGetAll_SuccessFromES(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.IRepository)
	mockAuthorRepo := new(authorMocks.IRepository)
	mockMQ := new(mqMocks.IRabbitMQClient)
	mockES := new(esMocks.ElasticsearchClient)
	mockCache := new(cacheMock.CacheRepository)

	req := structs.RequestSearchArticle{
		Keyword: "golang",
	}

	expectedResp := &structs.ResponseGetArticle{
		Total: 1,
		Article: []structs.Article{
			structs.Article{
				Title: "Sample",
				Body: "Isi",
			},
		},
	}
	
	sampleResEs :=  structs.ArticleESResponse{
		Hits: structs.Hits{
			Total: structs.Total{
				Value: 1,
			},
			Hits: []structs.Hit{
				structs.Hit{
					Source: structs.Source{
						ArticleEs: structs.ArticleEs{
							Title: "Sample",
							Body: "Isi",
						},
					},
				},
			},
		},
		Shards: structs.Shards{
			Total:      1,
			Successful: 1,
			Skipped:    0,
			Failed:     0,
		},
	}
	

	// Step 1: Cache miss
	mockCache.On("Get", ctx, req, mock.Anything).Return(errors.New("cache miss"))

	// Mocking Elasticsearch
	mockES.On("SearchInElasticsearch", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			dest := args.Get(3).(*structs.ArticleESResponse)
			*dest = sampleResEs
		}).Return(nil)

	// Mocking Cache
	mockCache.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("cache miss"))
	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	u := &ArticleUsecase{
		repo:     mockRepo,
		authorRepo: mockAuthorRepo,
		mqClient: mockMQ,
		esClient: mockES,
		conf:    &config.Config{
			ArticleIndex: "articles",
		},
		cache:    mockCache,
	}

	result, err := u.GetAll(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp.Total, result.Total)
	assert.Equal(t, expectedResp.Article[0].Title, result.Article[0].Title)

	mockCache.AssertExpectations(t)
	mockES.AssertExpectations(t)
}
