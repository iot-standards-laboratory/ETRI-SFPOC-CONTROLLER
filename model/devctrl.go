package model

import (
	"errors"
	"etri-sfpoc-controller/devmanager"
	"sync"
)

var devCtrls = map[string]devmanager.DeviceControllerI{}
var mutex sync.Mutex

func AddDeviceController(dname string, ctrl devmanager.DeviceControllerI) {
	mutex.Lock()
	defer mutex.Unlock()
	devCtrls[dname] = ctrl
}

func GetDeviceController(dname string) (devmanager.DeviceControllerI, error) {
	ctrl, ok := devCtrls[dname]
	if !ok {
		return nil, errors.New("does not exist error")
	}

	return ctrl, nil
}

func RemoveDeviceController(dname string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(devCtrls, dname)
}
