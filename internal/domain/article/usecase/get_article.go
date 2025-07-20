package usecase

import (
	"context"

	"github.com/dika22/news-service/package/structs"
)

func (u *ArticleUsecase) GetAll(ctx context.Context, req structs.RequestSearchArticle) (structs.ResponseGetArticle, error) {
	dest := structs.ResponseGetArticle{}
	if err := u.cache.Get(ctx, req, &dest); err == nil {
		return structs.ResponseGetArticle{}, nil
	}
	if dest.Total == 0 {
		respES := structs.ArticleESResponse{}
		query := req.NewQuerySearchArticle()
		if err := u.esClient.SearchInElasticsearch(ctx, u.conf.ArticleIndex, query, &respES); err != nil{
			return structs.ResponseGetArticle{}, err
		}

		dest = respES.NewResponseGetArticle()
		if err := u.cache.Set(ctx, req, &dest); err != nil {
			return structs.ResponseGetArticle{}, err
		}
		return dest, nil
	}
	return  dest, nil
}