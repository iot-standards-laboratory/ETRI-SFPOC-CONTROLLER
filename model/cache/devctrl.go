package cache

import (
	"errors"
	"etri-sfpoc-controller/devmanager"
	"sync"
)

var devCtrls = map[string]devmanager.DeviceControllerI{}
var ctrlMutex sync.Mutex

func AddDeviceController(dname string, ctrl devmanager.DeviceControllerI) {

	ctrlMutex.Lock()
	defer ctrlMutex.Unlock()
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
	ctrlMutex.Lock()
	defer ctrlMutex.Unlock()
	delete(devCtrls, dname)
}
