package usecase

import (
	"context"

	"github.com/dika22/news-service/package/structs"
)

func (u *ArticleUsecase) Create(ctx context.Context, req *structs.RequestCreateArticle) error {
	article := req.NewArticle()
	_, err := u.repo.Store(ctx, article)
	if err != nil {
		return err
	}
	return nil
}