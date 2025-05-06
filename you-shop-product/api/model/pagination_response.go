package model

type PaginationResponse struct {
	Page            int `json:"page"`
	PageSize        int `json:"page_size"`
	TotalItemsCount int `json:"total_items"`
	Items           any `json:"items"`
}

func NewPaginationResponse(page, pageSize, totalItems int, items any) *PaginationResponse {
	return &PaginationResponse{
		Page:            page,
		PageSize:        pageSize,
		TotalItemsCount: totalItems,
		Items:           items,
	}
}
