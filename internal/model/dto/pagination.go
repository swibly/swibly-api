package dto

type Pagination[T any] struct {
	Data []*T `json:"data"`

	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`
	CurrentPage  int `json:"current_page"`
	NextPage     int `json:"next_page"`
	PreviousPage int `json:"previous_page"`
}
