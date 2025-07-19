package structs

type PayloadMessageArticle struct {
	Articles Articles `json:"articles"`
	Authors Authors  `json:"authors"`
}

func (p PayloadMessageArticle) NewArticle() ArticlesElasticDocument {
	return ArticlesElasticDocument{
		Article: Article{
			ID:    p.Articles.ID,
			Title: p.Articles.Title,
			Body:  p.Articles.Body,
			Author: Author{
				ID:   p.Authors.ID,
				Name: p.Authors.Name,
			},
		},
		CreatedAt: p.Articles.CreatedAt,
	}
}