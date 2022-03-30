package notifier

type EventBox struct {
	INotiManager
}

var Box *EventBox

func init() {
	Box = &EventBox{NewNotiManager()}
}
