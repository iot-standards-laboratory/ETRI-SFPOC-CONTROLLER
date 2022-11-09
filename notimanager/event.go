package notimanager

type IEvent interface {
	Token() uint64
	Message() interface{} // Href
}

// Default implementation of IEvent
type SimpleEvent struct {
	_message interface{}
	_token   uint64
}

func NewSimpleEvent(_token uint64, _message interface{}) IEvent {
	return &SimpleEvent{
		_token:   _token,
		_message: _message,
	}
}

func (e *SimpleEvent) Token() uint64 {
	return e._token
}

func (e *SimpleEvent) Message() interface{} {
	return e._message
}
