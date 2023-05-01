package cachestorage

import (
	"bytes"
	"encoding/json"
	"errors"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/devmanager"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var devCtrls = map[uint64]devmanager.DeviceControllerI{}
var ctrlMutex sync.Mutex

func AddDeviceController(ctrl devmanager.DeviceControllerI) {

	ctrlMutex.Lock()
	defer ctrlMutex.Unlock()
	devCtrls[ctrl.Key()] = ctrl

	log.Println("added ctrl:", devCtrls)
	err := postController(ctrl)
	if err != nil {
		fmt.Println(err)
	}
}

func GetDeviceController(key uint64) (devmanager.DeviceControllerI, error) {
	ctrl, ok := devCtrls[key]
	if !ok {
		return nil, errors.New("does not exist error")
	}

	return ctrl, nil
}

func GetDeviceControllers() map[uint64]devmanager.DeviceControllerI {
	return devCtrls
}

func RemoveDeviceController(key uint64) {
	ctrlMutex.Lock()
	defer ctrlMutex.Unlock()

	err := deleteController(devCtrls[key])
	if err != nil {
		fmt.Println(err)
	}
	delete(devCtrls, key)

	log.Println("removed ctrl:", devCtrls)
}

func postController(e devmanager.DeviceControllerI) error {
	id, ok := config.Params["id"]
	if !ok {
		return errors.New("invalid id error")
	}
	edgeAddr, ok := config.Params["edgeAddress"]
	if !ok {
		return errors.New("invalid edgeAddr error")
	}

	payload := map[string]interface{}{
		"name":         e.Name(),
		"id":           fmt.Sprintf("%d", e.Key()),
		"agent_id":     id,
		"key":          fmt.Sprintf("%s/%d", id, e.Key()),
		"service_name": e.ServiceName(),
		"service_id":   e.ServiceID(),
	}

	bytesBuf := &bytes.Buffer{}
	enc := json.NewEncoder(bytesBuf)
	err := enc.Encode(payload)
	if err != nil {
		return err
	}

	resp, _ := http.Post(fmt.Sprintf("http://%s/api/v2/ctrls", edgeAddr), "application/json", bytesBuf)
	if resp.StatusCode != 200 {

		return errors.New("update failed error")
	}
	return nil
}

func deleteController(e devmanager.DeviceControllerI) error {
	id, ok := config.Params["id"]
	if !ok {
		return errors.New("invalid id error")
	}
	edgeAddr, ok := config.Params["edgeAddress"]
	if !ok {
		return errors.New("invalid edgeAddr error")
	}

	payload := map[string]interface{}{
		"name":         e.Name(),
		"id":           fmt.Sprintf("%d", e.Key()),
		"agent_id":     id,
		"key":          fmt.Sprintf("%s/%d", id, e.Key()),
		"service_name": e.ServiceName(),
		"service_id":   e.ServiceID(),
	}

	bytesBuf := &bytes.Buffer{}
	enc := json.NewEncoder(bytesBuf)
	err := enc.Encode(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("http://%s/api/v2/ctrls", edgeAddr),
		bytesBuf,
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("update failed error")
	}

	return nil
}
