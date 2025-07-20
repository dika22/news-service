package structs


type RequestCreateArticle struct {
	AuthorID int64  `json:"author_id" validate:"required"`
	Title    string `json:"title" validate:"required"`
	Body     string `json:"body" validate:"required"`
}