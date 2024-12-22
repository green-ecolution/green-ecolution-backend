package entities

type EventType string

type Event interface {
	Type() EventType
}

const (
	EventTypeUpdateTree        EventType = "update tree"
	EventTypeUpdateTreeCluster           = "update tree cluster"
)

type BasicEvent struct {
	eventType EventType
}

func (e BasicEvent) Type() EventType {
	return e.eventType
}

type EventUpdateTree struct {
	BasicEvent
	Old Tree
	New Tree
}

func NewEventUpdateTree(old, new Tree) EventUpdateTree {
	return EventUpdateTree{
		BasicEvent: BasicEvent{eventType: EventTypeUpdateTree},
		Old:        old,
		New:        new,
	}
}

type EventUpdateTreeCluster struct {
	BasicEvent
	Old TreeCluster
	New TreeCluster
}

func NewEventUpdateTreeCluster(old, new TreeCluster) EventUpdateTreeCluster {
	return EventUpdateTreeCluster{
		BasicEvent: BasicEvent{eventType: EventTypeUpdateTreeCluster},
		Old:        old,
		New:        new,
	}
}
