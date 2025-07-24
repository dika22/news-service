package usecase

import (
	"context"

	"github.com/dika22/news-service/internal/constant"
	"github.com/dika22/news-service/package/structs"
	"github.com/spf13/cast"
)


func (u *ArticleUsecase) GetByID(ctx context.Context, id int64) (structs.Article, error) {
	dest := structs.Article{}
	res, err, _ := u.group.Do(cast.ToString(id), func() (interface{}, error) {
		res, err := u.repo.GetByID(ctx, id)
		if err != nil {
			return structs.Article{}, err
			
		}
		return  res, nil
    })
	if err != nil {
		return structs.Article{}, err
	}
	
	art := res.(structs.Articles)
	author, err := u.authorRepo.GetByID(ctx, art.AuthorID)
	if err != nil {
		return structs.Article{}, err
	}

	dest.ID = art.ID
	dest.Title = art.Title
	dest.Body = art.Body
	dest.Status = constant.ArticleStatus[cast.ToInt(art.Status)]
	dest.Author = structs.Author{
		ID:   author.ID,
		Name: author.Name,
	}
	return  dest, nil
}