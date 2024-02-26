package response

type BasePagination struct {
	Total        int64 `json:"total"`
	CurrentPage  int   `json:"currentPage"`
	PrePageCount int   `json:"prePageCount"`
}

type Test struct {
}