package entities

type Pagination struct {
	Total       int64  `json:"total_records"`
	CurrentPage int32  `json:"current_page"`
	TotalPages  int32  `json:"total_pages"`
	NextPage    *int32 `json:"next_page"`
	PrevPage    *int32 `json:"prev_page"`
} // @Name Pagination

func CalculatePagination(totalCount, limit, page int32) (int32, *int32, *int32){
	totalPages := (totalCount + limit - 1) / limit
	var nextPage, prevPage *int32
	if page < totalPages {
		next := int32(page + 1)
		nextPage = &next
	}
	if page > 1 {
		prev := int32(page - 1)
		prevPage = &prev
	}

	return totalPages, nextPage, prevPage
}
