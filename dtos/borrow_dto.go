package dtos

type GetAllBorrowedDto struct {
	Skip  *int `form:"skip" binding:"required"`
	Limit *int `form:"limit" binding:"required"`
}

type BorrowDto struct {
	BookId     string `json:"bookId" bson:"bookId" binding:"required"`
	ReturnDate string `json:"returnDate" bson:"returnDate" binding:"required"`
}

type ApproveBookBorrowDto struct {
	Status int `json:"status" binding:"required,min=1,max=2"`
}

type ReturnBookBorrow struct {
}

type BorrowedHistoryDto struct {
	Skip  *int `form:"skip" binding:"required"`
	Limit int  `form:"limit" binding:"required"`
}
