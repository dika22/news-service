package usecase

import (
	"context"

	"github.com/dika22/news-service/internal/domain/article/repository"
	authorRepo "github.com/dika22/news-service/internal/domain/author/repository"
	"golang.org/x/sync/singleflight"

	"github.com/dika22/news-service/package/config"
	"github.com/dika22/news-service/package/connection/elasticsearch"
	rabbitmq "github.com/dika22/news-service/package/rabbit-mq"
	"github.com/dika22/news-service/package/structs"

	repoCache "github.com/dika22/news-service/internal/domain/article/repository/cache"
)

type IArticle interface{
	GetAll(ctx context.Context, req structs.RequestSearchArticle) (structs.ResponseGetArticle, error)
	Create(ctx context.Context, req *structs.RequestCreateArticle) error
	UpdatePublishArticle(ctx context.Context, req *structs.RequestUpdatePublishArticle) error
	GetByID(ctx context.Context, id int64) (structs.Article, error)
}

type ArticleUsecase struct{
	repo     repository.IRepository
	authorRepo authorRepo.IRepository
	mqClient rabbitmq.IRabbitMQClient
	esClient elasticsearch.ElasticsearchClient
	conf     *config.Config
	cache repoCache.CacheRepository
	group singleflight.Group
}


func NewsUsecase(repo repository.IRepository, 
		authorRepo authorRepo.IRepository,
		mqClient rabbitmq.IRabbitMQClient, 
		esClient elasticsearch.ElasticsearchClient,
		conf *config.Config,
		cache repoCache.CacheRepository) IArticle  {
	return &ArticleUsecase{
		repo: repo,
		authorRepo: authorRepo,
		mqClient: mqClient,
		esClient: esClient,
		conf: conf,
		cache: cache,
	}
}
