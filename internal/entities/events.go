package entities

type EventType string

type Event interface {
	Type() EventType
}

const (
	EventTypeUpdateTree        EventType = "update tree"
	EventTypeCreateTree                  = "create tree"
	EventTypeDeleteTree                  = "delete tree"
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
	Prev Tree
	New  Tree
}

func NewEventUpdateTree(prev, new Tree) EventUpdateTree {
	return EventUpdateTree{
		BasicEvent: BasicEvent{eventType: EventTypeUpdateTree},
		Prev:       prev,
		New:        new,
	}
}

type EventCreateTree struct {
	BasicEvent
	New Tree
}

func NewEventCreateTree(new Tree) EventCreateTree {
	return EventCreateTree{
		BasicEvent: BasicEvent{eventType: EventTypeCreateTree},
		New:        new,
	}
}

type EventDeleteTree struct {
	BasicEvent
	Prev Tree
}

func NewEventDeleteTree(prev Tree) EventDeleteTree {
	return EventDeleteTree{
		BasicEvent: BasicEvent{eventType: EventTypeDeleteTree},
		Prev:       prev,
	}
}

type EventUpdateTreeCluster struct {
	BasicEvent
	Prev TreeCluster
	New  TreeCluster
}

func NewEventUpdateTreeCluster(prev, new TreeCluster) EventUpdateTreeCluster {
	return EventUpdateTreeCluster{
		BasicEvent: BasicEvent{eventType: EventTypeUpdateTreeCluster},
		Prev:       prev,
		New:        new,
	}
}
