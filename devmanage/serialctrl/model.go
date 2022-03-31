package serialctrl

import "io"

type _device struct {
	UUID      string
	IfaceName string
	Sname     string
	Iface     io.ReadWriter
	states    map[string]interface{}
}

type RecvEvent struct {
	Params map[string]interface{}
}
