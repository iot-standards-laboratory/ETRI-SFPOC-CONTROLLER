package notifier

type IEvent interface {
	Token() string
	Title() string
	Body() interface{} // Href
}

// Default implementation of IEvent
type StatusChangedEvent struct {
	_title string
	_body  interface{}
	_token string
}

func (e *StatusChangedEvent) Token() string {
	return e._token
}

func (e *StatusChangedEvent) Title() string {
	return e._title
}

func (e *StatusChangedEvent) Body() interface{} {
	return e._body
}

func NewEvent(_title string, _body interface{}, _token string) IEvent {
	return &StatusChangedEvent{
		_title,
		_body,
		_token,
	}
}
