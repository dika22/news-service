package structs

import "time"

type Articles struct {
	ID       	int64      `json:"id"`
	AuthorID 	int64      `json:"author_id"`
	Title    	string     `json:"title"`
	Body     	string     `json:"body"`
	CreatedAt 	time.Time  `json:"created_at"`
	UpdatedAt 	time.Time  `json:"updated_at"`
}


func (p RequestCreateArticle) NewArticle() Articles {
	return Articles{
		AuthorID:   p.AuthorID,
		Title:      p.Title,
		Body:       p.Body,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}