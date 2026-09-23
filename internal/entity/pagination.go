package entity

// PaginationMeta contains pagination metadata for list endpoints
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalData   int64 `json:"total_data"`
	TotalPages  int   `json:"total_pages"`
}
