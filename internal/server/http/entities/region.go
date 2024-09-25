package entities

type RegionResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
} // @Name Region

type RegionListResponse struct {
	Regions    []*RegionResponse `json:"regions"`
	Pagination Pagination        `json:"pagination"`
} // @Name RegionList
