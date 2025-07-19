package structs

import "time"


type ArticlesElasticDocument struct {
	Article Article `json:"article"`
	CreatedAt time.Time `json:"created_at"`
}

type Article struct {
	ID       	int64      `json:"id"`
	Title    	string     `json:"title"`
	Body     	string     `json:"body"`
	Author   	Author     `json:"author"`
}

type Author struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
}

