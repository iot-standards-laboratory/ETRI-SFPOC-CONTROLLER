package model

import (
	"errors"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/devmanager"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

var DefaultDB DBHandlerI

func init() {
	var err error
	DefaultDB, err = NewSqliteHandler("dump.db")
	if err != nil {
		panic(err)
	}
}

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

func (s *_DBHandler) GetSID(sname string) (string, error) {
	sid, ok := s.sidCache[sname]
	if !ok {
		req, err := http.NewRequest("GET",
			fmt.Sprintf("http://%s/%s", config.Params["serverAddr"], "api/v1/svcs"),
			nil,
		)

		req.Header.Set("sname", sname)

		if err != nil {
			return "", err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		} else if resp.ContentLength == 0 {
			return "", errors.New("not exist service")
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		sid = string(b)
	}

	return sid, nil
}
