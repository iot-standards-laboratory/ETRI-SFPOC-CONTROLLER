package notifier

import (
	"github.com/gorilla/websocket"
)

type _WebsocketSubscriber struct {
	_id       string
	_token    string
	_type     int
	_complete chan<- int
	_conn     *websocket.Conn
	Msg       string
	// _h     func(msg string)
}

func (rs *_WebsocketSubscriber) String() string {
	return "id: " + rs._id
}
func (rs *_WebsocketSubscriber) ID() string {
	return rs._id
}

func (rs *_WebsocketSubscriber) Token() string {
	return rs._token
}

func (rs *_WebsocketSubscriber) Handle(e IEvent) {
	payload := map[string]interface{}{"key": e.Title()}

	if rs._conn.WriteJSON(payload) != nil {
		rs._complete <- 200
	}
}

func (rs *_WebsocketSubscriber) Type() int {
	return rs._type
}

func NewWebsocketSubscriber(_id, _token string, _type int, _complete chan<- int, _conn *websocket.Conn) ISubscriber {
	return &_WebsocketSubscriber{_id: _id, _token: _token, _type: _type, _complete: _complete, _conn: _conn, Msg: _id}

}
