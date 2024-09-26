package entities

type RegionResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
} // @Name Region

type RegionListResponse struct {
	Regions    []*RegionResponse `json:"regions"`
	Pagination Pagination        `json:"pagination"`
} // @Name RegionList
