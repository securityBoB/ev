package event

type Event struct {
	eventType  string
	eventData  interface{}
	dispatcher Dispatcher
}

func (e *Event) TypeName() string {
	return e.eventType
}

func (e *Event) Data() interface{} {
	return e.eventData
}

func (e *Event) Dispatcher() Dispatcher {
	return e.dispatcher
}

func (e *Event) SetDispatcher(dispatcher Dispatcher) {
	e.dispatcher = dispatcher
}

func NewEvent(eventType string, eventData interface{}) Event {

	return Event{eventType: eventType, eventData: eventData, dispatcher: NewDispatcher()}
}
