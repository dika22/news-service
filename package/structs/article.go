package structs

import (
	"time"

	"github.com/dika22/news-service/internal/constant"
)

type Articles struct {
	ID       	int64      `json:"id"`
	AuthorID 	int64      `json:"author_id"`
	Title    	string     `json:"title"`
	Body     	string     `json:"body"`
	Status      interface{}       `json:"status"` // 1 draft, 2 published 3 deleted
	CreatedAt 	time.Time  `json:"created_at"`
	UpdatedAt 	time.Time  `json:"updated_at"`
}


func (p RequestCreateArticle) NewArticle() Articles {
	return Articles{
		AuthorID:   p.AuthorID,
		Title:      p.Title,
		Body:       p.Body,
		Status:     constant.Drafted,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}