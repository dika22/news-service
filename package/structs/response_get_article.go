package structs

import "github.com/dika22/news-service/internal/constant"

type ResponseGetArticle struct {
	Page    int `json:"page"`
	Limit   int `json:"limit"`
	Total   int `json:"total"`
	Article []Article `json:"articles"`
}


func (res ArticleESResponse) NewResponseGetArticle() ResponseGetArticle {
	articles := []Article{}
	for _, hit := range res.Hits.Hits {
		article := Article{
			ID:     hit.Source.ArticleEs.ID,
			Title:  hit.Source.ArticleEs.Title,
			Body:   hit.Source.ArticleEs.Body,
			Status: constant.ArticleStatus[hit.Source.ArticleEs.Status],
			Author: Author{	
				ID:   hit.Source.ArticleEs.AuthorEs.ID,
				Name: hit.Source.ArticleEs.AuthorEs.Name,
			},
		}
		articles = append(articles, article)
	}
	return ResponseGetArticle{
		Total:   res.Hits.Total.Value,
		Article: articles,
	}
}