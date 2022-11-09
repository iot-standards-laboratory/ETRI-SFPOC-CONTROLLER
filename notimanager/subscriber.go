package notimanager

type ISubscriber interface {
	Handle(e IEvent)
	Type() int
	ID() string
	Token() uint64
}

// default implementation of ISubscriber
type callbackSubscriber struct {
	_id    string
	_token uint64
	_type  int
	_h     func(msg interface{})
}

func (rs *callbackSubscriber) ID() string {
	return rs._id
}

func (rs *callbackSubscriber) Token() uint64 {
	return rs._token
}

func (rs *callbackSubscriber) Handle(e IEvent) {
	rs._h(e.Message())
}

func (rs *callbackSubscriber) Type() int {
	return rs._type
}

func NewCallbackSubscriber(_id string, _token uint64, _type int, _h func(interface{})) ISubscriber {
	return &callbackSubscriber{_id: _id, _token: _token, _type: _type, _h: _h}
}

type chanSubscriber struct {
	_id      string
	_token   uint64
	_type    int
	_channel chan<- IEvent
}

func (rs *chanSubscriber) ID() string {
	return rs._id
}

func (rs *chanSubscriber) Token() uint64 {
	return rs._token
}

func (rs *chanSubscriber) Type() int {
	return rs._type
}

func (rs *chanSubscriber) Handle(e IEvent) {
	rs._channel <- e
}

func NewChanSubscriber(_id string, _token uint64, _type int, _channel chan<- IEvent) ISubscriber {
	return &chanSubscriber{_id: _id, _token: _token, _type: _type, _channel: _channel}
}

type webSocketSubscriber struct {
	_id       string
	_token    uint64
	_type     int
	_complete chan<- int
}

func (ws *webSocketSubscriber) String() string {
	return "id: " + ws._id
}

func (ws *webSocketSubscriber) ID() string {
	return ws._id
}

func (ws *webSocketSubscriber) Token() uint64 {
	return ws._token
}

func (ws *webSocketSubscriber) Handle(e IEvent) {
	payload := e.Message()
	_ = payload
}

func (ws *webSocketSubscriber) Type() int {
	return ws._type
}

func NewWebSocketSubscriber(_id string, _token uint64, _type int, _complete chan<- int) ISubscriber {
	return &webSocketSubscriber{
		_id,
		_token,
		_type,
		_complete,
	}
}
