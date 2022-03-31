package serialctrl

type Event interface {
	Params() map[string]interface{}
	Key() interface{}
}
type EventHandler interface {
	Handle(e Event)
}

type Receiver interface {
	onRecv(key interface{}, params map[string]interface{})
}

type Syncronizer interface {
	onSync(key interface{}, params map[string]interface{})
}
