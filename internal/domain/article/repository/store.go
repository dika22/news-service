package repository

import (
	"context"

	"github.com/dika22/news-service/package/structs"
)

func (r ArticleRepository) Store(ctx context.Context, payload structs.Articles) (int64, error)  {
	var id int64
	err := r.db.QueryRow(`INSERT INTO articles (title, body, author_id) VALUES ($1, $2, $3) RETURNING id`, payload.Title, payload.Body, payload.AuthorID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}