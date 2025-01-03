package utils

import (
	"errors"
	"log/slog"
	"slices"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

var mapWateringStatus = map[int]entities.WateringStatus{
	0: entities.WateringStatusGood,
	1: entities.WateringStatusModerate,
	2: entities.WateringStatusBad,
}

func mapKpaRange(centibar, lower, higher int) int {
	if centibar < lower {
		return 0
	} else if centibar < higher {
		return 1
	} else {
		return 2
	}
}

func CheckAndSortWatermarks(w []entities.Watermark) (w30, w60, w90 entities.Watermark, err error) {
	watermarks := slices.SortedFunc(slices.Values(w), func(a, b entities.Watermark) int {
		return a.Depth - b.Depth
	})

	if len(watermarks) != 3 || watermarks[0].Depth != 30 || watermarks[1].Depth != 60 || watermarks[2].Depth != 90 {
		err = errors.New("sensor data watermarks are malformed")
		return
	}

	w30, w60, w90 = watermarks[0], watermarks[1], watermarks[2]
	return
}

// CalculateWateringStatus determines the watering status of a plant based on its planting year and sensor watermarks.
//
// Parameters:
//   - plantingYear: The year the plant was planted.
//   - watermarks: A slice of entities.Watermark containing sensor readings at different depths.
//
// Returns:
//   - entities.WateringStatus: The calculated watering status based on the plant's lifetime and sensor watermarks.
//
// Behavior:
//  1. Calculates the plant's lifetime in years based on the current year.
//  2. Validates the watermarks to ensure exactly three readings exist at depths of 30cm, 60cm, and 90cm.
//     If validation fails, logs an error and returns `WateringStatusUnknown`.
//  3. Based on the tree's lifetime, applies specific centibar ranges to the watermarks at different depths:
//     - Year 1: Centibar ranges for 30cm, 60cm, and 90cm are 25–33.
//     - Year 2: Centibar ranges for 30cm are 62–81, while 60cm and 90cm remain 25–33.
//     - Year 3: Centibar ranges for 30cm are 1585–infinity, while 60cm and 90cm are 80–infinity.
//  4. Maps the centibar values to a status (green, yellow, or red) and determines the final watering status
//     based on the most severe status.
//
// Errors:
//   - If the watermarks are malformed or do not contain exactly three readings at the required depths, logs an error
//     with details and returns `WateringStatusUnknown`.
//
// Example:
//
//	plantingYear := 2020
//	watermarks := []entities.Watermark{
//	    {Depth: 30, Centibar: 28},
//	    {Depth: 60, Centibar: 30},
//	    {Depth: 90, Centibar: 35},
//	}
//
//	status := CalculateWateringStatus(plantingYear, watermarks)
//	fmt.Printf("Watering Status: %v\n", status)
//
// Notes:
//   - The function assumes that watermarks are provided as a slice, where each entry represents a sensor reading at a specific depth.
//   - Any changes to the mapping of centibar ranges or tree lifetime logic should be reflected here.
func CalculateWateringStatus(plantingYear int32, watermarks []entities.Watermark) entities.WateringStatus {
	/*
		Tree 1st year:
		30cm: <25kPA: green; 25-32kPA orange; >32kPA red
		60cm: <25kPA: green; 25-32kPA orange; >32kPA red
		90cm: <25kPA: green; 25-32kPA orange; >32kPA red

		Tree 2nd year:
		30cm: <62kPA: green; 62-80kPA orange; >80kPA red
		60cm: <25kPA: green; 25-32kPA orange; >32kPA red
		90cm: <25kPA: green; 25-32kPA orange; >32kPA red

		Tree 3rd year:
		30cm: <1585kPa: green;
		60cm: <80kPA: green; >80kPA red
		90cm: <80kPA: green; >80kPA red
	*/
	currentYear := int32(time.Now().Year())
	treeLifetime := currentYear - plantingYear
	w30, w60, w90, err := CheckAndSortWatermarks(watermarks)
	if err != nil {
		slog.Error("sensor data watermarks are malformed", "watermarks", watermarks)
		return entities.WateringStatusUnknown
	}

	statusList := make([]int, 3)
	switch treeLifetime {
	case 0:
		fallthrough
	case 1:
		statusList[0] = mapKpaRange(w30.Centibar, 25, 33)
		statusList[1] = mapKpaRange(w60.Centibar, 25, 33)
		statusList[2] = mapKpaRange(w90.Centibar, 25, 33)
	case 2:
		statusList[0] = mapKpaRange(w30.Centibar, 62, 81)
		statusList[1] = mapKpaRange(w60.Centibar, 25, 33)
		statusList[2] = mapKpaRange(w90.Centibar, 25, 33)
	case 3:
		statusList[0] = mapKpaRange(w30.Centibar, 1586, -1) // -1 means no orange
		statusList[1] = mapKpaRange(w60.Centibar, 80, -1)
		statusList[2] = mapKpaRange(w90.Centibar, 80, -1)
	default:
		return entities.WateringStatusUnknown
	}

	slices.Sort(statusList)
	return mapWateringStatus[statusList[2]]
}
