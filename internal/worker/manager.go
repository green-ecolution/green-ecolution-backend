package worker

import (
	"context"
	"errors"
	"sync"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

var (
	UnknownEventTypeErr = errors.New("unsupported event type")
	NotSubscribedErr    = errors.New("not subscribed")
)

type Subscriber interface {
	HandleEvent(ctx context.Context, event entities.Event) error
	EventType() entities.EventType
}

type EventManager struct {
	eventCh    chan entities.Event
	subscriber map[entities.EventType]map[int]chan<- entities.Event
	nextID     int
	eventTypes map[entities.EventType]struct{}
	rwMutex    sync.RWMutex
}

func NewEventManager(eventTypes ...entities.EventType) *EventManager {
	eventTypeMap := make(map[entities.EventType]struct{})
	for _, eventType := range eventTypes {
		eventTypeMap[eventType] = struct{}{}
	}

	subscriber := make(map[entities.EventType]map[int]chan<- entities.Event)
	for eventType := range eventTypeMap {
		subscriber[eventType] = make(map[int]chan<- entities.Event)
	}

	return &EventManager{
		eventCh:    make(chan entities.Event, 1),
		subscriber: subscriber,
		nextID:     0,
		eventTypes: eventTypeMap,
	}
}

func (e *EventManager) Publish(event entities.Event) error {
	if _, ok := e.eventTypes[event.Type()]; !ok {
		return UnknownEventTypeErr
	}

	e.eventCh <- event

	return nil
}

func (e *EventManager) Subscribe(eventType entities.EventType) (int, <-chan entities.Event, error) {
	if _, ok := e.eventTypes[eventType]; !ok {
		return -1, nil, UnknownEventTypeErr
	}

	e.rwMutex.Lock()
	defer e.rwMutex.Unlock()

	ch := make(chan entities.Event)
	id := e.nextID
	e.subscriber[eventType][id] = ch
	e.nextID++

	return id, ch, nil
}

func (e *EventManager) Unsubscribe(eventType entities.EventType, id int) error {
	if _, ok := e.eventTypes[eventType]; !ok {
		return UnknownEventTypeErr
	}

	e.rwMutex.Lock()
	defer e.rwMutex.Unlock()

	return e.unsubscribe(eventType, id)
}

func (e *EventManager) unsubscribe(eventType entities.EventType, id int) error {
	ch, ok := e.subscriber[eventType][id]
	if !ok {
		return NotSubscribedErr
	}

	close(ch)
	delete(e.subscriber[eventType], id)

	return nil
}

func (e *EventManager) cleanup() {
	e.rwMutex.Lock()
	defer e.rwMutex.Unlock()

	for eventType, subscriber := range e.subscriber {
		for id := range subscriber {
			_ = e.unsubscribe(eventType, id)
		}

		delete(e.eventTypes, eventType)
	}
}

func (e *EventManager) Stop() {
	close(e.eventCh)
}

// Blocking, caller has to run this in goroutine
func (e *EventManager) Run(ctx context.Context) {
	defer e.cleanup()
	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-e.eventCh:
			if !ok {
				return
			}

			// TODO: Maybe unsubscribe if can't write
			e.rwMutex.RLock()
			for _, ch := range e.subscriber[v.Type()] {
				select {
				case <-ctx.Done():
					e.rwMutex.RUnlock()
					return
				case ch <- v:
				}
			}
			e.rwMutex.RUnlock()
		}
	}
}

func (em *EventManager) RunSubscription(ctx context.Context, sub Subscriber) error {
	id, ch, err := em.Subscribe(sub.EventType())
	if err != nil {
		return err
	}
	defer em.Unsubscribe(sub.EventType(), id)
	for {
		select {
		case <-ctx.Done():
			return nil
		case v, ok := <-ch:
			if !ok {
				return nil
			}
			if err := sub.HandleEvent(ctx, v); err != nil {
				return err
			}
		}
	}
}
