package repository

import (
	"context"
	"database/sql"

	"github.com/dika22/news-service/package/structs"
)

type AuthorRepository struct {
	db *sql.DB
}

type IRepository interface {
	GetByID(ctx context.Context, id int64) (structs.Authors, error)
}

func NewAuthorRepository(db *sql.DB) IRepository {
	return AuthorRepository{
		db: db,
	}
}