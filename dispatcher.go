package event

import "sync"

type Dispatcher interface {
	Register(eventType string, handler Handler)
	AddListener(eventType string, listener *Listener)
	RemoveListener(eventType string, listener *Listener) bool
	HasListener(eventType string) bool
	DispatchEvent(event Event) bool
	Trigger(eventType string, params interface{}) bool
}

type eventDispatcher struct {
	savers map[string]*Saver
	sync.RWMutex
}

func NewDispatcher() Dispatcher {
	return &eventDispatcher{savers: make(map[string]*Saver)}
}

func (dispatcher *eventDispatcher) Register(eventType string, handler Handler) {
	dispatcher.AddListener(eventType, NewListener(handler))
}

func (dispatcher *eventDispatcher) AddListener(eventType string, listener *Listener) {
	dispatcher.Lock()
	defer dispatcher.Unlock()

	saver, ok := dispatcher.savers[eventType]
	if !ok {
		saver = &Saver{Type: eventType, Listeners: []*Listener{listener}}
		dispatcher.savers[eventType] = saver
		return
	}

	saver.Lock()
	defer saver.Unlock()

	saver.Listeners = append(saver.Listeners, listener)
}

func (dispatcher *eventDispatcher) RemoveListener(eventType string, listener *Listener) bool {
	dispatcher.Lock()
	defer dispatcher.Unlock()

	saver, ok := dispatcher.savers[eventType]
	if !ok {
		return false
	}

	saver.Lock()
	defer saver.Unlock()

	for i, l := range saver.Listeners {
		if l == listener {
			saver.Listeners = append(saver.Listeners[:i], saver.Listeners[i+1:]...)
			return true
		}
	}
	return false
}

func (dispatcher *eventDispatcher) HasListener(eventType string) bool {
	dispatcher.RLock()
	defer dispatcher.RUnlock()

	_, ok := dispatcher.savers[eventType]
	return ok
}

func (dispatcher *eventDispatcher) DispatchEvent(event Event) bool {
	dispatcher.RLock()
	defer dispatcher.RUnlock()

	saver, ok := dispatcher.savers[event.TypeName()]
	if !ok {
		return false
	}

	saver.RLock()
	defer saver.RUnlock()

	for _, listener := range saver.Listeners {
		event.SetDispatcher(dispatcher)
		go listener.Handler(event)
	}

	return true
}

func (dispatcher *eventDispatcher) Trigger(eventType string, params interface{}) bool {
	event := NewEvent(eventType, params)
	return dispatcher.DispatchEvent(event)
}
