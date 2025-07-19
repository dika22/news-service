package usecase_test

import (
	"context"
	"testing"

	"news-service/internal/domain/article/repository/mocks"
	"news-service/internal/domain/article/usecase"
	authorMocks "news-service/internal/domain/author/repository/mocks"
	"news-service/package/config"
	"news-service/package/structs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mqMocks "news-service/package/rabbit-mq/mocks"
)

func TestCreate_Success(t *testing.T) {
	// Mock Dependencies
	mockRepo := new(mocks.IRepository)
	mockAuthorRepo := new(authorMocks.IRepository)
	mockMQ := new(mqMocks.IRabbitMQClient)

	usecase := usecase.NewsUsecase(
		mockRepo,
		mockAuthorRepo,
		mockMQ,
		nil, // esClient
		&config.Config{ArticleQueue: "article.created"},
		nil,
	)

	req := &structs.RequestCreateArticle{
		Title:    "Judul",
		Body:     "Konten",
		AuthorID: 1,
	}

	article := req.NewArticle()
	article.ID = 1

	mockRepo.On("Store", mock.Anything, mock.Anything).Return(int64(1), nil)
	mockAuthorRepo.On("GetByID", mock.Anything, int64(1)).Return(structs.Authors{
		ID:   1,
		Name: "Adhika",
	}, nil)
	mockMQ.On("Publish", "article.created", mock.Anything).Return(nil)

	err := usecase.Create(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockAuthorRepo.AssertExpectations(t)
	mockMQ.AssertExpectations(t)
}

