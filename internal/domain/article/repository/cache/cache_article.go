package repository

import (
	"context"
	"fmt"
	"news-service/package/structs"
	"time"
)

func (r ArticleCacheRepository) Get(ctx context.Context, req structs.RequestSearchArticle, dest *structs.ResponseGetArticle) error {
	key := fmt.Sprintf("article:%v:%v", req.Page, req.Limit)
	return r.cache.Get(ctx, key, dest)
}

func (r ArticleCacheRepository) Set(ctx context.Context, req structs.RequestSearchArticle, dest *structs.ResponseGetArticle) error  {
	key := fmt.Sprintf("article:%v:%v", req.Page, req.Limit)
	return r.cache.Set(ctx, key, dest, 180 * time.Second)
}