package structs


type RequestCreateArticle struct {
	AuthorID int64 `json:"author_id"`
	Title string `json:"title"`
	Body string `json:"body"`
}