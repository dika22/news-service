package structs

type ArticleESResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards Shards `json:"_shards"`
	Hits Hits `json:"hits"`
}

type Shards   struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type Hits struct {
	Total Total `json:"total"`
	MaxScore interface{} `json:"max_score"`
	Hits []Hit `json:"hits"`
}

type Hit struct {
	Index string `json:"_index"`
	Type string `json:"_type"`
	ID string `json:"_id"`
	Score float64 `json:"_score"`
	Source Source `json:"_source"`
}

type Source struct {
	ArticleEs ArticleEs `json:"article"`
	CreatedAt string `json:"created_at"`
}

type ArticleEs struct {
	ID       int64        `json:"id"`
	Title    string     `json:"title"`
	Body     string     `json:"body"`
	Status   int `json:"status"`
	AuthorEs AuthorEs `json:"author"`
}

type AuthorEs struct {
	ID   int64    `json:"id"`
	Name string `json:"name"`
}

type Total struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}
