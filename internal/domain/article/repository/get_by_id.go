package repository

import (
	"context"

	"github.com/dika22/news-service/package/structs"
)


func (r ArticleRepository) GetByID(ctx context.Context, id int64) (structs.Articles, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM articles WHERE id = $1", id)
	if err != nil {
		return structs.Articles{}, err
	}
	article := structs.Articles{}
	for rows.Next() {
		err := rows.Scan(&article.ID, &article.AuthorID, &article.Body, &article.Title, &article.CreatedAt, &article.UpdatedAt, &article.Status)
		if err != nil {
			return structs.Articles{}, err
		}
	}
	return article, nil
}