package response

type BasePagination struct {
	Count       int64 `json:"count"`
	CurrentPage int   `json:"current_page"`
	PageCount   int   `json:"page_count"`
}

type Test struct {
}
