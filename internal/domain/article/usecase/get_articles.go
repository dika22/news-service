package usecase

import (
	"context"
	"fmt"

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
		resp, err, _ := u.group.Do(fmt.Sprintf("articles:%v:%v", req.Page, req.Limit), func() (interface{}, error) {
			err := u.esClient.SearchInElasticsearch(ctx, u.conf.ArticleIndex, query, &respES) 
			if err != nil{
				return structs.ResponseGetArticle{}, err
			}
			return respES, nil
		})
		respES = resp.(structs.ArticleESResponse)
		if err != nil {
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
	return  dest, nil
}