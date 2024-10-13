package fileimport

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/omniscale/go-proj/v2"
)

type GeoTransformer struct {
	from        *proj.Proj
	to          *proj.Proj
	transformer proj.Transformer
}

func NewGeoTransformer(from, to int) (*GeoTransformer, error) {
	fromProj, err := proj.NewEPSG(from)
	if err != nil {
		return nil, err
	}

	toProj, err := proj.NewEPSG(to)
	if err != nil {
		return nil, err
	}

	transformer, err := proj.NewEPSGTransformer(from, to)
	if err != nil {
		return nil, err
	}

	return &GeoTransformer{
		from:        fromProj,
		to:          toProj,
		transformer: transformer,
	}, nil
}

func (g *GeoTransformer) Transform(x, y float64) (lat, lng float64, err error) {
	points := []proj.Coord{
		proj.XY(x, y),
	}

	if err := g.transformer.Transform(points); err != nil {
		return 0, 0, err
	}

	return points[0].X, points[0].Y, nil
}

type GeoPoint struct {
	X float64
	Y float64
}

func (g *GeoTransformer) TransformBatch(points []GeoPoint) ([]GeoPoint, error) {
	coords := utils.Map(points, func(p GeoPoint) proj.Coord {
		return proj.XY(p.X, p.Y)
	})

	if err := g.transformer.Transform(coords); err != nil {
		return nil, err
	}

	return utils.Map(coords, func(c proj.Coord) GeoPoint {
		return GeoPoint{X: c.X, Y: c.Y}
	}), nil
}
