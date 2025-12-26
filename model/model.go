package model

// model data

type Book struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Stock int    `json:"stock"`
}

type Member struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Quota int    `json:"quota"`
}

type BorrowRequest struct {
	MemberID int `json:"member_id"`
	BookID   int `json:"book_id"`
}

// custom respon error

type CustomErrorResponse struct {
	Message        string `json:"message"`
	ZiyadErrorCode string `json:"ziyad_error_code"`
	TraceID        string `json:"trace_id"`
}
