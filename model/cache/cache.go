package cache

import (
	"encoding/json"
	"errors"
	"etri-sfpoc-controller/common"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/devmanager"
	"etri-sfpoc-controller/model"
	"fmt"
	"sync"
)

var devCtrls = map[string]devmanager.DeviceControllerI{}
var ctrlMutex sync.Mutex

type devicemeta struct {
	Did   string `json:"did"`
	Dname string `json:"dname"`
	Cid   string `json:"cid"`
}

// services to registered device
var svcs = map[string][]*devicemeta{} // svcIds[sname] = meta of device list
var svcMutex sync.Mutex

// when AddService is called : Device registration is done
// always called with AddDeviceController
func AddSvc(did, dname, sname string) error {
	// start mutex
	svcMutex.Lock()
	defer svcMutex.Unlock()

	// get list for sname
	svcs[sname] = append(svcs[sname], &devicemeta{
		Did:   did,
		Dname: dname,
		Cid:   config.Params["cid"].(string),
	})

	return nil
}

func GetSvcList() map[string][]*devicemeta {
	fmt.Println("GetSvcList(): ", svcs)
	return svcs
}

func GetDevicesOnSvc(sname string) []*model.Device {
	list, ok := svcs[sname]
	if !ok {
		return nil
	}

	devList := make([]*model.Device, len(list))
	for i, e := range list {
		device := &model.Device{
			DID:   e.Did,
			DName: e.Dname,
			CID:   e.Cid,
		}

		devList[i] = device
	}

	return devList
}

func RemoveDeviceFromSvc(did string) error {

	// get sname for device
	device, err := model.DefaultDB.GetDevice(did)
	if err != nil {
		return err
	}
	sname := device.SName
	// when remove service
	svcMutex.Lock()
	defer svcMutex.Unlock()

	// get list for sname
	list, ok := svcs[sname]
	if !ok {
		return errors.New("not exist service")
	}

	// remove dname in the list
	for i, e := range list {
		if e.Did == did {
			list[i] = list[len(svcs)-1]
			if len(list)-1 == 0 {
				delete(svcs, sname)
				removeSvcId(sname)
			} else {
				svcs[sname] = list[:len(list)-1]
			}
		}
	}

	return nil
	// remove sname entity when the service
}

func subcribeSvc(sid string) {
	cid := config.Params["cid"].(string)
	go common.Subscribe(
		fmt.Sprintf("svc/%s/%s", sid, "push/v1/"),
		cid,
		func(payload []byte) {
			cmdJson := map[string]interface{}{}
			err := json.Unmarshal(payload, &cmdJson)
			if err != nil {
				return
			}

			key, ok := cmdJson["key"].(string)
			if !ok {
				return
			}

			if key == "control" {
				value, ok := cmdJson["value"].(map[string]interface{})
				if !ok {
					return
				}
				did, ok := value["did"].(string)
				if !ok {
					return
				}

				dev, err := model.DefaultDB.GetDevice(did)
				if err != nil {
					return
				}

				ctrl, err := GetDeviceController(dev.DName)
				if err != nil {
					panic(err)
				}

				status, ok := value["status"].(map[string]interface{})
				fmt.Println(status)
				if !ok {
					return
				}
				ctrl.Sync(status)
			}
		})
}

var svcIds = map[string]string{} // svcIds[sname] = sid
var svcIdMutex sync.Mutex

//
func AddSvcId(sname, sid string) error {
	// lock mutex
	svcIdMutex.Lock()
	defer svcIdMutex.Unlock()

	// check whether the service already exist
	_, ok := svcIds[sname]
	if ok {
		return errors.New("already exist service")
	}

	// add service
	svcIds[sname] = sid

	// start subscribing the service
	subcribeSvc(sid)
	return nil
}

func GetSvcId(sname string) (string, bool) {
	svcIdMutex.Lock()
	defer svcIdMutex.Unlock()
	id, ok := svcIds[sname]
	return id, ok
}

func GetSvcIds() map[string]string {
	// unsubscribe service
	// Todo
	return svcIds
}

// if all devices managed by service with the sname is removed, this method is called
func removeSvcId(sname string) {
	svcIdMutex.Lock()
	defer svcIdMutex.Unlock()
	delete(svcIds, sname)
}

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

// func (s *_DBHandler) GetSID(sname string) (string, error) {
// 	sid, ok := s.sidCache[sname]
// 	if !ok {
// 		req, err := http.NewRequest("GET",
// 			fmt.Sprintf("http://%s/%s", config.Params["serverAddr"], "api/v1/svcs"),
// 			nil,
// 		)

// 		req.Header.Set("sname", sname)

// 		if err != nil {
// 			return "", err
// 		}

// 		resp, err := http.DefaultClient.Do(req)
// 		if err != nil {
// 			return "", err
// 		} else if resp.ContentLength == 0 {
// 			return "", errors.New("not exist service")
// 		}

// 		b, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			return "", err
// 		}

// 		sid = string(b)
// 	}

// 	return sid, nil
// }
