package entities

type WateringStatus string

const (
	WateringStatusGood     WateringStatus = "good"
	WateringStatusModerate WateringStatus = "moderate"
	WateringStatusBad      WateringStatus = "bad"
	WateringStatusUnknown  WateringStatus = "unknown"
)
