package cache

import (
	"context"
	"encoding/json"
	"errors"
	"etri-sfpoc-controller/commonutils"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/model"
	"fmt"
	"strings"
	"sync"
)

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

func subcribeSvc(sname, sid string) context.CancelFunc {
	cid := config.Params["cid"].(string)
	ctx, cancel := context.WithCancel(context.Background())
	go commonutils.Subscribe(
		ctx,
		fmt.Sprintf("svc/%s/%s", sid, "push/v1/"),
		cid,
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
		},
		func() {
			removeSvcId(sname)
		},
	)

	return cancel
}

var cancels = map[string]context.CancelFunc{}
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
	cancels[sname] = subcribeSvc(sname, sid)
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

	cancel, ok := cancels[sname]
	if ok {
		cancel()
		delete(cancels, sname)
	}
}

func GetSvcUrls(sname, path string) (string, error) {

	if path[0] != '/' {
		return "", errors.New("path should start '/'")
	}

	var sid string
	var ok bool
	if strings.Compare(config.Params["mode"].(string), string(config.STANDALONE)) == 0 {

		sid, ok = config.Params["sid"].(string)
		if !ok {
			return "", errors.New("sid is blank")
		} else if strings.Compare(sid, "blank") == 0 {
			return "", errors.New("sid is blank")
		}
	} else {
		sid, ok = GetSvcId(sname)
		if !ok {
			return "", errors.New("sid is blank")
		}
	}

	return fmt.Sprintf("http://%s/svc/%s%s", config.Params["serverAddr"], sid, path), nil
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
