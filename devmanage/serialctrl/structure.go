package serialctrl

type EventStruct struct {
	Event
	params map[string]interface{}
	key    interface{}
}

func (es *EventStruct) Params() map[string]interface{} {
	return es.params
}

func (es *EventStruct) Key() interface{} {
	return es.key
}

type EventHandlerStruct struct {
	EventHandler
	HandleFunc func(e Event)
}

func (ehs *EventHandlerStruct) Handle(e Event) {
	ehs.HandleFunc(e)
}

func NewEventHandler(h func(e Event)) *EventHandlerStruct {
	return &EventHandlerStruct{HandleFunc: h}
}
