package repository

import (
	"context"
	"fmt"

	"github.com/dika22/news-service/package/structs"
)

func (r ArticleRepository) Update(ctx context.Context, payload *structs.RequestUpdatePublishArticle) (int64, error)  {
	var id int64
	err := r.db.QueryRow(`UPDATE articles SET title = $1, body = $2, author_id = $3 , status = $4 WHERE id = $5 RETURNING id`, 
		payload.Title, payload.Body, payload.AuthorID, payload.ID, payload.Status).
		Scan(&id)
	if err != nil {
		fmt.Println("error update", err)
		return 0, err
	}

	fmt.Println("debug id", id)
	return id, nil
	
}