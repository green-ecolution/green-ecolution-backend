package entities

import (
	"strings"

	"github.com/pkg/errors"
)

type WateringStatus string

const (
	WateringStatusGood        WateringStatus = "good"
	WateringStatusModerate    WateringStatus = "moderate"
	WateringStatusBad         WateringStatus = "bad"
	WateringStatusJustWatered WateringStatus = "just watered"
	WateringStatusUnknown     WateringStatus = "unknown"
)

func ParseWateringStatus(status string) ([]WateringStatus, error) {
	if strings.Contains(status, ",") {
		parts := strings.Split(status, ",")
		var statuses []WateringStatus

		for _, part := range parts {
			parsedStatus, err := parseSingleWateringStatus(strings.TrimSpace(part))
			if err != nil {
				return nil, err
			}
			statuses = append(statuses, parsedStatus)
		}
		return statuses, nil
	}

	parsedStatus, err := parseSingleWateringStatus(status)
	if err != nil {
		return nil, err
	}
	return []WateringStatus{parsedStatus}, nil
}

func parseSingleWateringStatus(status string) (WateringStatus, error) {
	switch status {
	case string(WateringStatusGood):
		return WateringStatusGood, nil
	case string(WateringStatusModerate):
		return WateringStatusModerate, nil
	case string(WateringStatusBad):
		return WateringStatusBad, nil
	case string(WateringStatusUnknown):
		return WateringStatusUnknown, nil
	default:
		return "", errors.New("invalid watering status: " + status)
	}
}
