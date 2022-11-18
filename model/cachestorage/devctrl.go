package cachestorage

import (
	"errors"
	"etri-sfpoc-controller/devmanager"
	"log"
	"sync"
)

var devCtrls = map[uint64]devmanager.DeviceControllerI{}
var ctrlMutex sync.Mutex

func AddDeviceController(ctrl devmanager.DeviceControllerI) {

	ctrlMutex.Lock()
	defer ctrlMutex.Unlock()
	devCtrls[ctrl.Key()] = ctrl
	log.Println("added ctrl:", devCtrls)
}

func GetDeviceController(key uint64) (devmanager.DeviceControllerI, error) {
	ctrl, ok := devCtrls[key]
	if !ok {
		return nil, errors.New("does not exist error")
	}

	return ctrl, nil
}

func RemoveDeviceController(key uint64) {
	ctrlMutex.Lock()
	defer ctrlMutex.Unlock()
	delete(devCtrls, key)
	log.Println("removed ctrl:", devCtrls)
}
