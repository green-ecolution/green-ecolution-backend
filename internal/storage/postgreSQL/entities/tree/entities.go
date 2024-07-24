package tree

import "github.com/google/uuid"

type TreeLocationEntity struct {
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Address        string  `json:"address"`
	AdditionalInfo string  `json:"additional_info"`
}

type TreeEntity struct {
	ID       uuid.UUID          `json:"id"`
	Species  string             `json:"species"`
	TreeNum  int                `json:"tree_num"`
	Age      int                `json:"age"`
	Location TreeLocationEntity `json:"location"`
}
