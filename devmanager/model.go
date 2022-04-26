package devmanager

// import (
// 	"errors"
// 	"fmt"
// 	"io"
// 	"sync"
// )

// var mutex sync.Mutex
// var ifaces = map[string]interface{}{}
// var uuids = map[interface{}]string{}

// func addDeviceInterface(uuid string, iface io.ReadWriter) error {
// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	oldIface, ok := ifaces[uuid]
// 	if ok {
// 		if iface == oldIface {
// 			return errors.New("already exist uuid")
// 		}
// 		delete(uuids, oldIface)
// 		fmt.Println("change interface!!")
// 	}

// 	ifaces[uuid] = iface
// 	uuids[iface] = uuid

// 	return nil
// }

// func removeDeviceInterface(uuid string, iface io.ReadWriter) error {
// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	delete(uuids, iface)
// 	delete(ifaces, uuid)

// 	return nil
// }
