package entities

type Pagination struct {
	Total       int64  `json:"total_records"`
	CurrentPage int32  `json:"current_page"`
	TotalPages  int32  `json:"total_pages"`
	NextPage    *int32 `json:"next_page"`
	PrevPage    *int32 `json:"prev_page"`
} // @Name Pagination

func CalculatePagination(totalCount, limit, page int32) (totalPages int32, nextPage, prevPage *int32) {
	if limit <= 0 {
		limit = 1
	}

	if page < totalPages {
		next := page + 1
		nextPage = &next
	}

	if page > 1 {
		prev := page - 1
		prevPage = &prev
	}

	if page == totalPages {
		nextPage = nil
	}
	
	if page == 1 {
		prevPage = nil
	}

	return (totalCount + limit - 1) / limit, nextPage, prevPage
}
