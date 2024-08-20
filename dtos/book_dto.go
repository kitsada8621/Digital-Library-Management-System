package dtos

type GetBooksDto struct {
	Skip       *int64   `json:"skip" binding:"required"`
	Limit      int64    `json:"limit" binding:"required"`
	Search     *string  `json:"search" binding:"required"` // Book title/ author
	Categories []string `json:"categories" binding:"required"`
}

type BookDto struct {
	CategoryId      string `json:"categoryId" binding:"required"`
	Author          string `json:"author" binding:"required,max=255"`
	BookTitle       string `json:"bookTitle" binding:"required,max=255"`
	BookDescription string `json:"bookDesc" binding:"required"`
}
