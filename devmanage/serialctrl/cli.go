package serialctrl

import (
	"context"
	"etri-sfpoc-controller/devmanage/serialctrl/puserial"
	"log"
	"sync"

	"github.com/rjeczalik/notify"
)

var ctx context.Context
var cancel context.CancelFunc
var ch_discover chan notify.EventInfo
var _managerObj *_manager
var registerHandleFunc func(e Event)
var removeHandleFunc func(e Event)

func init() {
	ctx, cancel = context.WithCancel(context.Background())
	ch_discover = make(chan notify.EventInfo)
	devWithUUID := map[string]*_device{}
	devWithIface := map[interface{}]*_device{}
	chanForSync := map[string]chan map[string]interface{}{}

	_managerObj = &_manager{
		devicesWithUUID:  devWithUUID,
		devicesWithIface: devWithIface,
		chanForSync:      chanForSync,
		SyncListener:     &SyncHandler{devices: devWithUUID, chanForSync: chanForSync, mutex: &sync.Mutex{}, states: map[string]map[string]interface{}{}},
		RecvListener:     &RecvHandler{devices: devWithIface, chanForSync: chanForSync},
	}

	registerHandleFunc = nil
	removeHandleFunc = nil
}

func AddRegisterHandleFunc(h func(e Event)) {
	registerHandleFunc = h
}

func AddRemoveHandleFunc(h func(e Event)) {
	removeHandleFunc = h
}

func Close() {
	cancel()
}

func Sync(key string, param map[string]interface{}) {
	_managerObj.onSync(key, param)
}

func AddRecvListener(h EventHandler) {
	_managerObj.addRecvListener(h)
}

// func SetDevicePropsToSync(uuid string, propsToSync []string) error {
// 	device, ok := _managerObj.devicesWithUUID[uuid]
// 	if !ok {
// 		return errors.New("device not found")
// 	}
// 	device.propsToSync = propsToSync
// 	return nil
// }

func Run() error {
	ifaces, err := puserial.InitDevice()
	if err != nil {
		return err
	}

	for _, e := range ifaces {
		go _managerObj.onAdded(e)
	}

	go puserial.WatchNewDevice(ctx, ch_discover)

	for {
		e, ok := <-ch_discover
		if !ok {
			log.Println("manager exit")
			return nil
		}
		switch e.Event() {
		case notify.Create:
			go _managerObj.onAdded(e.Path())
			// case notify.Remove:
			// 	log.Println("USB Disconnected!!")
			// 	_managerObj.onAdded(e.Path())
		}

	}
}
