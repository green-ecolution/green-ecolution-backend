package fileimport

import "github.com/omniscale/go-proj/v2"

type GeoTransformer struct {
	from        proj.Proj
	to          proj.Proj
	transformer proj.Transformer
}

func NewGeoTransformer(from, to proj.Proj) *GeoTransformer {
	return &GeoTransformer{
		from: from,
		to:   to,
	}
}
