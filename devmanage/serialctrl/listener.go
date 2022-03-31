package serialctrl

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sync"
)

type SyncHandler struct {
	devices     map[string]*_device
	mutex       *sync.Mutex
	chanForSync map[string]chan map[string]interface{}
	states      map[string]map[string]interface{}
}

func compareMap(src map[string]interface{}, dst map[string]interface{}) bool {

	for key, value := range src {
		if key == "code" {
			continue
		}

		if reflect.TypeOf(value).String() == "string" {
			dstV, ok := dst[key].(string)
			if !ok || value != dstV {
				return false
			}
			continue
		}

		srcV, ok := value.(float64)
		if !ok {
			srcV = float64(value.(int))
		}

		dstV, ok := dst[key].(float64)
		if !ok {
			dstV = float64(dst[key].(int))
		}

		if srcV != dstV {
			fmt.Println("diff ", key, "] ", srcV, " vs ", dstV)
			return false
		}
	}

	return true
}

func (sh *SyncHandler) Handle(e Event) {
	fmt.Println("Sync")
	device, ok := sh.devices[e.Key().(string)]
	// fmt.Println("sync] ", device.IfaceName)
	// fmt.Println(sh.devices)
	if !ok {
		return
	}

	params := e.Params()
	if params == nil {
		return
	}

	go func() {
		encoder := json.NewEncoder(device.Iface)
		origin := map[string]interface{}{}
		origin["code"] = 200
		for key, value := range params {
			origin[key] = value
		}

		sh.states[device.UUID] = origin
		// props := []string{"fan", "light", "servo"}

		err := encoder.Encode(origin)
		if err != nil {
			return
		}

		_, ok := sh.chanForSync[device.IfaceName]
		if ok {
			sh.mutex.Lock()
			close(sh.chanForSync[device.IfaceName])
			delete(sh.chanForSync, device.IfaceName)
			sh.mutex.Unlock()
		}

		chanForSync := make(chan map[string]interface{})
		sh.mutex.Lock()
		sh.chanForSync[device.IfaceName] = chanForSync
		sh.mutex.Unlock()

		for state := range chanForSync {
			if compareMap(sh.states[device.UUID], state) {
				sh.mutex.Lock()
				close(sh.chanForSync[device.IfaceName])
				delete(sh.chanForSync, device.IfaceName)
				sh.mutex.Unlock()
				return
			}
			log.Println("wrong: ", state)
			log.Println("resend: ", sh.states[device.UUID])
			err := encoder.Encode(sh.states[device.UUID])
			if err != nil {
				return
			}
		}
	}()
}

type RecvHandler struct {
	devices     map[interface{}]*_device
	chanForSync map[string]chan map[string]interface{}
	next        EventHandler
}

func (rh *RecvHandler) Handle(e Event) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic recover - ", r)
		}
	}()
	device, ok := rh.devices[e.Key()]
	if !ok {
		return
	}

	param := e.Params()
	code, _ := param["code"].(float64)
	if int(code) == 100 {
		return
	}

	device.states = param
	reportParam := map[string]interface{}{"state": device.states, "sname": device.Sname, "did": device.Did}
	channel, ok := rh.chanForSync[device.IfaceName]
	if ok {
		channel <- param
	}

	if rh.next != nil {
		rh.next.Handle(&EventStruct{key: e.Key(), params: reportParam})
	}

	// fmt.Println("recv] ", device.states)
}
