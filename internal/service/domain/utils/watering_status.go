package utils

import (
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

// CalculateWateringStatus determines the watering status of a plant based on its planting year and sensor data.
//
// Parameters:
//   - plantingYear: The year the plant was planted.
//   - data: A pointer to entities.SensorData, which contains watermarks and other sensor readings.
//
// Returns:
//   - entities.WateringStatus: The calculated watering status based on the plant's lifetime and sensor data.
//
// Behavior:
//  1. Calculates the plant's lifetime in years based on the current year.
//  2. Validates the watermarks from the sensor data to ensure exactly three readings exist at depths of 30cm, 60cm, and 90cm.
//     If validation fails, logs an error and returns `WateringStatusUnknown`.
//  3. Based on the tree's lifetime, applies specific centibar ranges to the watermarks at different depths:
//     - Year 1: Centibar ranges for 30cm, 60cm, and 90cm are 25–33.
//     - Year 2: Centibar ranges for 30cm are 62–81, while 60cm and 90cm remain 25–33.
//     - Year 3: Centibar ranges for 30cm are 1585–infinity, while 60cm and 90cm are 80–infinity.
//  4. Maps the centibar values to a status (green, orange, or red) and determines the final watering status
//     based on the most severe status.
//
// Errors:
//   - If the watermarks are malformed or incomplete, logs an error with details and returns `WateringStatusUnknown`.
//
// Example:
//
//	plantingYear := 2020
//	sensorData := &entities.SensorData{
//	    Data: entities.SensorReadings{
//	        Watermarks: map[int]entities.Watermark{
//	            30: {Depth: 30, Centibar: 28},
//	            60: {Depth: 60, Centibar: 30},
//	            90: {Depth: 90, Centibar: 35},
//	        },
//	    },
//	}
//
//	status := CalculateWateringStatus(plantingYear, sensorData)
//	fmt.Printf("Watering Status: %v\n", status)
//
// Notes:
//   - The function assumes a specific structure of sensor data where watermarks are mapped by depth.
//   - Any changes to the mapping of centibar ranges or tree lifetime logic should be reflected here.
func CalculateWateringStatus(plantingYear int32, data *entities.SensorData) entities.WateringStatus {
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
	treeLifetime := plantingYear - currentYear

	watermarks := slices.SortedFunc(slices.Values(data.Data.Watermarks), func(a, b entities.Watermark) int {
		return a.Depth - b.Depth
	})

	if len(watermarks) != 3 || watermarks[0].Depth != 30 || watermarks[1].Depth != 60 || watermarks[2].Depth != 90 {
		slog.Error("sensor data watermarks are malformed", "watermarks", watermarks)
		return entities.WateringStatusUnknown
	}

	watermark30, watermark60, watermark90 := watermarks[0], watermarks[1], watermarks[2]

	statusList := make([]int, 3)
	switch treeLifetime {
	case 1:
		statusList[0] = mapKpaRange(watermark30.Centibar, 25, 33)
		statusList[1] = mapKpaRange(watermark60.Centibar, 25, 33)
		statusList[2] = mapKpaRange(watermark90.Centibar, 25, 33)
	case 2:
		statusList[0] = mapKpaRange(watermark30.Centibar, 62, 81)
		statusList[1] = mapKpaRange(watermark60.Centibar, 25, 33)
		statusList[2] = mapKpaRange(watermark90.Centibar, 25, 33)
	case 3:
		statusList[0] = mapKpaRange(watermark30.Centibar, 1585, -1) // -1 means no orange
		statusList[1] = mapKpaRange(watermark60.Centibar, 80, -1)
		statusList[2] = mapKpaRange(watermark90.Centibar, 80, -1)
	default:
		return entities.WateringStatusUnknown
	}

	slices.Sort(statusList)
	return mapWateringStatus[statusList[2]]
}
