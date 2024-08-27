package tree

type TreeLocationResponse struct {
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Address        string  `json:"address"`
	AdditionalInfo string  `json:"additional_info"`
} // @Name TreeLocation

type TreeResponse struct {
	ID       int32                `json:"id"`
	Species  string               `json:"species"`
	TreeNum  int32                `json:"tree_num"`
	Age      int32                `json:"age"`
	Location TreeLocationResponse `json:"location"`
} // @Name Tree
