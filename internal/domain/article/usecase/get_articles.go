package usecase

import (
	"context"

	"github.com/dika22/news-service/package/structs"
)

func (u *ArticleUsecase) GetAll(ctx context.Context, req structs.RequestSearchArticle) (structs.ResponseGetArticle, error) {
	dest := structs.ResponseGetArticle{}
	if dest.Total > 1 {
		if err := u.cache.Get(ctx, req, &dest); err == nil {
			return structs.ResponseGetArticle{}, nil
		}
		return  dest, nil
	}
	respES := structs.ArticleESResponse{}
	query := req.NewQuerySearchArticle()
	err := u.esClient.SearchInElasticsearch(ctx, u.conf.ArticleIndex, query, &respES) 
	if err != nil{
		return structs.ResponseGetArticle{}, err
	}
	dest = respES.NewResponseGetArticle()
	dest.Limit = query["size"].(int)
	dest.Page  = query["from"].(int) / dest.Limit + 1
	if err := u.cache.Set(ctx, req, &dest); err != nil {
		return structs.ResponseGetArticle{}, err
	}
	return dest, nil
}