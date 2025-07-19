package repository

import (
	"context"
	"log"
	"news-service/package/structs"
)

func (r ArticleRepository) GetAll(ctx context.Context) ([]*structs.Articles, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM articles")
	if err != nil {
		return nil, err
	}
	
	dest := []*structs.Articles{}
	for rows.Next() {
		article := structs.Articles{}
		err := rows.Scan(&article.ID, &article.AuthorID, &article.Body, &article.Title, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			log.Fatal(err)
		}
		dest = append(dest, &article)
	}
	// Cek jika ada kesalahan dalam iterasi
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
		
	return  dest, err
}