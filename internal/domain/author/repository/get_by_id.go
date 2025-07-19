package repository

import (
	"context"
	"news-service/package/structs"
)

func (r AuthorRepository) GetByID(ctx context.Context, id int64) (structs.Authors, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM authors WHERE id = $1", id)
	if err != nil {
		return structs.Authors{}, err
	}
	author := structs.Authors{}
	for rows.Next() {
		err := rows.Scan(&author.ID, &author.Name, &author.CreatedAt, &author.UpdatedAt)
		if err != nil {
			return structs.Authors{}, err
		}
	}
	return author, nil
}