package repository

import (
	"context"
	"database/sql"

	"github.com/dika22/news-service/package/connection/cache"
	"github.com/dika22/news-service/package/structs"
)

type IRepository interface{
	GetAll(ctx context.Context) ([]*structs.Articles, error)
	Store(ctx context.Context, payload structs.Articles) (int64, error)
	Update(ctx context.Context, payload *structs.RequestUpdatePublishArticle) (int64, error)
	GetByID(ctx context.Context, id int64) (structs.Articles, error)
}


type ArticleRepository struct{
	db *sql.DB
	cache cache.Cache
}

func NewsRepository(db *sql.DB, cache cache.Cache) IRepository {
	return ArticleRepository{
		db: db,
		cache: cache,
	}
}