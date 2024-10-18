package request

import "context"

type BaseRequest struct {
	Ctx       context.Context `json:"ctx"`
	RequestId string          `json:"request_id"`
	Language  string          `json:"language"`
}

type BaseTokenRequest struct {
}

type BasePagination struct {
	CurrentPage int `json:"current_page" validate:"required,min=1"`
	PageCount   int `json:"page_count" validate:"required,max=50"`
}
