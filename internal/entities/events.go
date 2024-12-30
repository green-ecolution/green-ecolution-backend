package entities

type EventType string

type Event interface {
	Type() EventType
}

const (
	EventTypeUpdateTree        EventType = "update tree"
	EventTypeCreateTree        EventType = "create tree"
	EventTypeDeleteTree        EventType = "delete tree"
	EventTypeUpdateTreeCluster EventType = "update tree cluster"
	EventTypeNewSensorData     EventType = "receive sensor data"
)

type BasicEvent struct {
	eventType EventType
}

func (e BasicEvent) Type() EventType {
	return e.eventType
}

type EventUpdateTree struct {
	BasicEvent
	Prev *Tree
	New  *Tree
}

func NewEventUpdateTree(prev, newTree *Tree) EventUpdateTree {
	return EventUpdateTree{
		BasicEvent: BasicEvent{eventType: EventTypeUpdateTree},
		Prev:       prev,
		New:        newTree,
	}
}

type EventCreateTree struct {
	BasicEvent
	New *Tree
}

func NewEventCreateTree(newTree *Tree) EventCreateTree {
	return EventCreateTree{
		BasicEvent: BasicEvent{eventType: EventTypeCreateTree},
		New:        newTree,
	}
}

type EventDeleteTree struct {
	BasicEvent
	Prev *Tree
}

func NewEventDeleteTree(prev *Tree) EventDeleteTree {
	return EventDeleteTree{
		BasicEvent: BasicEvent{eventType: EventTypeDeleteTree},
		Prev:       prev,
	}
}

type EventUpdateTreeCluster struct {
	BasicEvent
	Prev *TreeCluster
	New  *TreeCluster
}

func NewEventUpdateTreeCluster(prev, newTc *TreeCluster) EventUpdateTreeCluster {
	return EventUpdateTreeCluster{
		BasicEvent: BasicEvent{eventType: EventTypeUpdateTreeCluster},
		Prev:       prev,
		New:        newTc,
	}
}

type EventNewSensorData struct {
	BasicEvent
	New *SensorData
}

func NewEventSensorData(newData *SensorData) EventNewSensorData {
	return EventNewSensorData{
		BasicEvent: BasicEvent{eventType: EventTypeNewSensorData},
		New:        newData,
	}
}
