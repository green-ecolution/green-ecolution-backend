package route

import (
	"log"
	"math"
)

type RouteService struct {
}

type Point struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
} // @Name Point

func NewRouteService() *RouteService {
	return &RouteService{}
}

func (s *RouteService) Ready() bool {
	return true
}

func (s *RouteService) DemoPoints() []Point {
	return []Point{
		{Name: "TBZ", Lat: 54.768822266460894, Lon: 9.434689073552093},
		{Name: "Point 0", Lat: 54.782792, Lon: 9.424908},
		{Name: "Point 1", Lat: 54.782792, Lon: 9.424908},
		{Name: "Point 2", Lat: 54.786913, Lon: 9.408921},
		{Name: "Point 3", Lat: 54.781590, Lon: 9.424873},
		{Name: "Point 4", Lat: 54.788817, Lon: 9.425888},
		{Name: "Point 5", Lat: 54.782141, Lon: 9.429827},
		{Name: "Point 6", Lat: 54.787557, Lon: 9.438296},
		{Name: "Point 7", Lat: 54.811338, Lon: 9.455175},
		{Name: "Point 8", Lat: 54.788306, Lon: 9.444110},
		{Name: "Point 9", Lat: 54.805982, Lon: 9.447570},
		{Name: "Point 10", Lat: 54.784699, Lon: 9.438025},
		{Name: "Point 11", Lat: 54.760139, Lon: 9.380937},
		{Name: "Point 12", Lat: 54.762046, Lon: 9.385276},
		{Name: "Point 13", Lat: 54.761091, Lon: 9.385719},
		{Name: "Point 14", Lat: 54.769905, Lon: 9.473510},
		{Name: "Point 15", Lat: 54.771202, Lon: 9.430948},
		{Name: "Point 16", Lat: 54.792941, Lon: 9.462763},
		{Name: "Point 17", Lat: 54.797287, Lon: 9.454632},
		{Name: "Point 18", Lat: 54.804966, Lon: 9.486204},
	}
}

func (s *RouteService) calculateDistance(p1 Point, p2 Point) float64 {
	radius := 6371.0 // Earth radius in km
	dLat := (p2.Lat - p1.Lat) * (math.Pi / 180.0)
	dLon := (p2.Lon - p1.Lon) * (math.Pi / 180.0)
	a := math.Pow(math.Sin(dLat/2), 2) +
		math.Cos(p1.Lat*(math.Pi/180.0))*
			math.Cos(p2.Lat*(math.Pi/180.0))*
			math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := radius * c
	log.Printf("Distance between %v and %v is %v km", p1, p2, distance)
	return distance
}

func (s *RouteService) Distance(points []Point) [][]float64 {
	n := len(points)
	distances := make([][]float64, n)
	for i := 0; i < n; i++ {
		distances[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			dist := s.calculateDistance(points[i], points[j])
			distances[i][j] = dist
		}
	}
	return distances
}

func (s *RouteService) NearestNeighbor(points []Point) ([]Point, float64) {
	n := len(points)
	visited := make([]bool, n)
	path := make([]int, n+1)
	totalDistance := 0.0
	distances := s.Distance(points)

	currentPoint := 0
	path[0] = currentPoint
	visited[0] = true

	for i := 0; i < n-1; i++ {
		nearest := 0
		nearestDistance := math.Inf(1)
		for point := 0; point < n; point++ {
			if visited[point] {
				continue
			}

			dist := distances[currentPoint][point]
			if dist < nearestDistance {
				nearestDistance = dist
				nearest = point
			}
		}

		currentPoint = nearest
		path[i+1] = currentPoint
		visited[currentPoint] = true
		totalDistance += nearestDistance
	}

	path[n] = 0
	totalDistance += distances[currentPoint][0]

	log.Printf("Path: %v", path)
	log.Printf("Total distance: %v", totalDistance)

	// Convert path to points
	var pathPoints []Point
	for _, p := range path {
		pathPoints = append(pathPoints, points[p])
	}

	return pathPoints, totalDistance
}
