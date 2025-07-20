package repository

import (
	"context"

	"github.com/dika22/news-service/package/connection/cache"
	"github.com/dika22/news-service/package/structs"
)

type ArticleCacheRepository struct{
	cache cache.Cache
}


type CacheRepository interface{
	Get(ctx context.Context, req structs.RequestSearchArticle, dest *structs.ResponseGetArticle) error
	Set(ctx context.Context, req structs.RequestSearchArticle, dest *structs.ResponseGetArticle) error
	DeleteArticleKeys(ctx context.Context) error
}


func NewCacheRepository(cache cache.Cache) CacheRepository {
	return ArticleCacheRepository{
		cache: cache,
	}
}