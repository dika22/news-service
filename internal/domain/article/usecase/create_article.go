package usecase

import (
	"context"
	"encoding/json"
	"log"
	"news-service/package/structs"
)

func (u *ArticleUsecase) Create(ctx context.Context, req *structs.RequestCreateArticle) error {
	article := req.NewArticle()
	id, err := u.repo.Store(ctx, article)
	if err != nil {
		return err
	}
	article.ID = id

	author, err := u.authorRepo.GetByID(ctx, article.AuthorID)
	if err != nil {
		return err
	}

	articleWithAuthor := structs.PayloadMessageArticle{
		Articles: article,
		Authors: author,
	}

	msg, errMarshal := json.Marshal(articleWithAuthor)
	if err != nil {
		log.Println(errMarshal)
		return nil
	}
	
	err = u.mqClient.Publish(u.conf.ArticleQueue, []byte(msg))
	if err != nil {
		return err
	}
	
	return nil
}