package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/dika22/news-service/package/structs"
)

func (r ArticleCacheRepository) Get(ctx context.Context, req structs.RequestSearchArticle, dest *structs.ResponseGetArticle) error {
	key := fmt.Sprintf("articles:%v:%v", req.Page, req.Limit)
	return r.cache.Get(ctx, key, dest)
}

func (r ArticleCacheRepository) Set(ctx context.Context, req structs.RequestSearchArticle, dest *structs.ResponseGetArticle) error  {
	key := fmt.Sprintf("articles:%v:%v", req.Page, req.Limit)
	return r.cache.Set(ctx, key, dest, 180 * time.Second)
}

func (r ArticleCacheRepository) DeleteArticleKeys(ctx context.Context) error {
	var cursor uint64
	var batchSize int64 = 100
	pattern := "articles:*"

	for {
		// Panggil Scan sesuai interface kamu
		keys, nextCursor, err := r.cache.Scan(ctx, int64(cursor), batchSize, pattern)
		if err != nil {
			return fmt.Errorf("failed to scan: %v", err)
		}

		if len(keys) > 0 {
			// Hapus key yang ditemukan
			if err := r.cache.Del(ctx, keys); err != nil {
				return fmt.Errorf("failed to delete keys: %v", err)
			}
			fmt.Printf("Deleted keys: %v\n", keys)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}
