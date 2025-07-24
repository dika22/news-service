package structs


type RequestUpdatePublishArticle struct {
	ID       int64 `json:"id"`
	AuthorID int64  `json:"author_id" validate:"required"`
	Title    string `json:"title" validate:"required"`
	Body     string `json:"body" validate:"required"`
	Status   int8 `json:"status" validate:"required"`
}