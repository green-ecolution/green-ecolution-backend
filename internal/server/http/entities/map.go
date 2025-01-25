package entities

type MapData struct {
	Center    		GeoJSONLocation 	`json:"center"`
	BoundSouthWest  GeoJSONLocation 	`json:"bounds_south_west"`
	BoundNorthEast 	GeoJSONLocation 	`json:"bounds_north_east"`
} // @Name MapData