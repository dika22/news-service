package repository

import (
	"context"
	"news-service/package/connection/cache"
	"news-service/package/structs"
)

type ArticleCacheRepository struct{
	cache cache.Cache
}


type CacheRepository interface{
	Get(ctx context.Context, req structs.RequestSearchArticle, dest *structs.ResponseGetArticle) error
	Set(ctx context.Context, req structs.RequestSearchArticle, dest *structs.ResponseGetArticle) error
}


func NewCacheRepository(cache cache.Cache) CacheRepository {
	return ArticleCacheRepository{
		cache: cache,
	}
}