package notifier

type ISubscriber interface {
	Handle(e IEvent)
	Type() int
	ID() string
	Token() string
}

// default implementation of ISubscriber
type CallbackSubscriber struct {
	_id    string
	_token string
	_type  int
	_h     func(e IEvent)
}

func (rs *CallbackSubscriber) ID() string {
	return rs._id
}

func (rs *CallbackSubscriber) Token() string {
	return rs._token
}

func (rs *CallbackSubscriber) Handle(e IEvent) {
	rs._h(e)
}

func (rs *CallbackSubscriber) Type() int {
	return rs._type
}

func NewCallbackSubscriber(_id, _token string, _type int, _h func(IEvent)) ISubscriber {
	return &CallbackSubscriber{_id: _id, _token: _token, _type: _type, _h: _h}
}

type ChanSubscriber struct {
	_id      string
	_token   string
	_type    int
	_channel chan<- IEvent
}

func (rs *ChanSubscriber) ID() string {
	return rs._id
}

func (rs *ChanSubscriber) Token() string {
	return rs._token
}

func (rs *ChanSubscriber) Type() int {
	return rs._type
}

func (rs *ChanSubscriber) Handle(e IEvent) {
	rs._channel <- e
}

func NewChanSubscriber(_id, _token string, _type int, _channel chan<- IEvent) ISubscriber {
	return &ChanSubscriber{_id: _id, _token: _token, _type: _type, _channel: _channel}
}
