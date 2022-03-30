package devmanage

import "etri-sfpoc-controller/notifier"

type EventBox struct {
	notifier.INotiManager
}

var eventBox *EventBox

func init() {
	eventBox = &EventBox{notifier.NewNotiManager()}
}
