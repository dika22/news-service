package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/dika22/news-service/internal/constant"
	"github.com/dika22/news-service/package/structs"
)


func (u *ArticleUsecase) UpdatePublishArticle(ctx context.Context, req *structs.RequestUpdatePublishArticle) error {
	article, err := u.repo.GetByID(ctx, req.ID)
	if err != nil {
		fmt.Println("error get by id", err, req.ID)
		return err
	}


	fmt.Println("debug article", article)
	if article.ID < 1 {
		return errors.New("article not found")
	}

	if article.Status == constant.Published {
		return errors.New("article is already published")
	}

	if req.ID != article.ID {
		return errors.New("article id not match")
		
	}

	id, err := u.repo.Update(ctx, req)
	if err != nil {
		fmt.Println("error updated", req)
		return err
	}

	article.ID = id
	author, err := u.authorRepo.GetByID(ctx, article.AuthorID)
	if err != nil {
		return err
	}

	articleWithAuthor := structs.PayloadMessageArticle{
		Articles: structs.Articles{
			ID:        req.ID,
			AuthorID:  req.AuthorID,
			Title:     req.Title,
			Body:      req.Body,
			Status:    req.Status,
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
		},
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