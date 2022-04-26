package devmanager

type Event interface {
	Params() map[string]interface{}
	Key() interface{}
}

type eventStruct struct {
	key    interface{}
	params map[string]interface{}
}

func (es *eventStruct) Params() map[string]interface{} {
	return es.params
}

func (es *eventStruct) Key() interface{} {
	return es.key
}

func NewEvent(params map[string]interface{}, key interface{}) Event {
	return &eventStruct{params: params, key: key}
}
