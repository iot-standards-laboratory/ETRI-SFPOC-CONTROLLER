package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/router"
	"etri-sfpoc-controller/statmgmt"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func register() (string, error) {
	// Controller 이름을 읽어옴
	payload := map[string]interface{}{}
	payload["cname"] = config.Params["cname"]
	fmt.Println("cname: ", payload["cname"])
	fmt.Println(config.Params["cname"])
	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Controller 등록 메시지 송신
	resp, err := http.Post(
		fmt.Sprintf("http://%s/%s", config.Params["serverAddr"], "api/v2/ctrls"),
		"application/json",
		bytes.NewReader(b),
	)

	if err != nil {
		return "", err
	}

	// 응답 메시지 수신
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	json.Unmarshal(b, &payload)

	// 등록 후 생성된 Controller ID 저장
	config.Set("cid", payload["cid"].(string))

	return payload["cid"].(string), nil
}

// return sid or error with record not found
// func querySvcID(sname string) (string, error) {

// 	if strings.Compare(config.Params["mode"].(string), string(config.STANDALONE)) == 0 {
// 		sid, _ := config.Params["sid"].(string)
// 		return sid, nil
// 	}

// 	// get svc id
// 	req, err := http.NewRequest(
// 		"GET",
// 		fmt.Sprintf("http://%s/%s", config.Params["serverAddr"], "api/v1/svcs"),
// 		nil,
// 	)

// 	if err != nil {
// 		return "", err
// 	}
// 	req.Header.Add("sname", sname)

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return "", err
// 	} else if resp.StatusCode != 200 {
// 		b, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			return "", err
// 		}
// 		return "", errors.New(string(b))
// 	}

// 	b, err := ioutil.ReadAll(resp.Body)
// 	return string(b), err
// }

// func deviceManagerSetup() {
// 	devmanager.AddOnDiscovered(func(port io.ReadWriter, sname, dname string) error {
// 		defer log.Println("exit onDiscovered()")
// 		// do register procedure
// 		//check already registered device
// 		cid := config.Params["cid"].(string)
// 		var did string
// 		var err error
// 		did, err = model.DefaultDB.GetDeviceID(dname)
// 		if err != nil {
// 			if strings.Compare(config.Params["mode"].(string), string(config.MANAGEDBYEDGE)) == 0 {
// 				fmt.Println("register device to edge")
// 				did, err = devmanager.RegisterDeviceToEdge(map[string]interface{}{
// 					"sname": sname,
// 					"dname": dname,
// 				})
// 				if err != nil {
// 					log.Println(err.Error())
// 					return err
// 				}
// 				err := devmanager.PostDeviceToEdge(did)
// 				if err != nil {
// 					log.Println(err.Error())
// 				}
// 			} else {
// 				_uuid, err := uuid.NewUUID()
// 				if err != nil {
// 					log.Println(err)
// 					return err
// 				}
// 				did = _uuid.String()
// 			}

// 			dev := &model.Device{
// 				DID:   did,
// 				DName: dname,
// 				SName: sname,
// 				CID:   cid,
// 			}

// 			model.DefaultDB.AddDevice(dev)
// 		} else {
// 			if strings.Compare(config.Params["mode"].(string), string(config.MANAGEDBYEDGE)) == 0 {
// 				fmt.Println("POSTDEVICETOEDGE")
// 				err := devmanager.PostDeviceToEdge(did)
// 				if err != nil {
// 					log.Println(err.Error())
// 				}
// 			}
// 		}

// 		// send request to server for registration of device
// 		okchan := make(chan error)
// 		defer close(okchan)
// 		ticker := time.NewTicker(10 * time.Second)
// 		defer ticker.Stop()

// 		_, err = port.Write([]byte(`{"code": 1, "token": "initial", "mode": 1}\n`))
// 		if err != nil {
// 			return err
// 		}

// 		go func() {
// 			reader := bufio.NewReader(port)
// 			rcvMsg := map[string]interface{}{}

// 			for {
// 				line, _, err := reader.ReadLine()
// 				if err != nil {
// 					log.Println("okchan : ", err)
// 					okchan <- err
// 					return
// 				}

// 				err = json.Unmarshal(line, &rcvMsg)
// 				if err != nil {
// 					continue
// 				}

// 				fmt.Println("code: ", rcvMsg["code"])
// 				if rcvMsg["code"].(float64)-1.0 < 0.0001 {
// 					// register success
// 					ctrl := makeDeviceController(port, did, dname, sname)
// 					cache.AddDeviceController(dname, ctrl)
// 					cache.AddSvc(did, dname, sname)
// 					sid, err := querySvcID(sname)
// 					if err == nil {
// 						// add service id to cache and subscribe service
// 						cache.AddSvcId(sname, sid)

// 						// register device to service if service is running
// 						dev := &model.Device{
// 							DID:   did,
// 							DName: dname,
// 							SName: sname,
// 							CID:   cid,
// 						}
// 						err = registerDeviceToService(sname, dev)
// 						if err != nil {
// 							log.Println(err)
// 						}
// 					} else {
// 						log.Println(err)
// 					}

// 					ctrl.Run()
// 					okchan <- nil
// 					return
// 				}
// 			}
// 		}()

// 		for {
// 			select {
// 			case <-ticker.C:
// 				fmt.Println("retransmission command to change mode as timeout")
// 				_, err = port.Write([]byte(`{"code": 1, "token": "initial", "mode": 1}\n`))
// 				if err != nil {
// 					return err
// 				}
// 			case err := <-okchan:
// 				if err == nil {
// 					return nil
// 				} else {
// 					log.Println(err)
// 				}
// 			}
// 		}
// 	})

// 	go devmanager.Watch()
// }

// func registerDeviceToService(sname string, dev *model.Device) error {

// 	// example of body
// 	// {
// 	// 	"did": "1ad73898-fe73-4053-8bbd-29a6ac6c9276",
// 	// 	"dname": "DEVICE-A-UUID",
// 	// 	"cid": "8971e21d-4800-47e4-b50e-691e28d2f701",
// 	// 	"sname": "devicemanagera"
// 	// }
// 	body, err := json.Marshal(dev)
// 	if err != nil {
// 		return err
// 	}

// 	svcUrl, err := cache.GetSvcUrls(sname, "/api/v1/devs")
// 	if err != nil {
// 		return err
// 	}
// 	log.Println("send device registration request to", svcUrl)
// 	// register device to service

// 	resp, err := http.Post(
// 		svcUrl,
// 		"application/json",
// 		bytes.NewReader(body),
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	respMsg, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}
// 	log.Println(respMsg)
// 	if resp.StatusCode > 300 {
// 		return errors.New(string(respMsg))
// 	} else {
// 		notifier.Box.Publish(
// 			notifier.NewEvent("connected", nil, sname),
// 		)
// 		return nil
// 	}
// }

// func makeDeviceController(port io.ReadWriter, did, dname, sname string) devmanager.DeviceControllerI {
// 	fmt.Println("makeDeviceController()")
// 	// model.AddDeviceController 에서 등록된 디바이스 목록에 해당 디바이스 추가할 것!!
// 	ctrl := devmanager.NewDeviceController(port, dname, did)

// 	msmtUrl := ""

// 	// when called device is registered on service!!
// 	subscriber := notifier.NewCallbackSubscriber(
// 		did,
// 		sname,
// 		notifier.SubtypeCont,
// 		func(i notifier.IEvent) {
// 			if strings.Compare(i.Title(), "connected") == 0 {

// 				svcUrl, err := cache.GetSvcUrls(sname, "/api/v1/status")
// 				if err != nil {
// 					log.Println(err.Error())
// 					return
// 				}

// 				log.Println("service is connected!! :", svcUrl)

// 				b, _ := json.Marshal(map[string]interface{}{
// 					"did":   did,
// 					"dname": dname,
// 					"cid":   config.Params["cid"],
// 				})
// 				// send request to "http://server/svc/{service_id}/api/v1/"
// 				req, err := http.NewRequest(
// 					"POST",
// 					svcUrl,
// 					bytes.NewReader(b),
// 				)

// 				if err != nil {
// 					log.Println(err)
// 					return
// 				}

// 				resp, err := http.DefaultClient.Do(req)
// 				if err != nil {
// 					log.Println(err)
// 					return
// 				} else if resp.StatusCode >= 300 {
// 					log.Println(err)
// 					return
// 				}

// 				b, err = ioutil.ReadAll(resp.Body)
// 				if err != nil {
// 					log.Println(err)
// 					return
// 				}

// 				msmtUrl = svcUrl + "/" + string(b)
// 			} else {
// 				log.Println("connection with service is disconnected")
// 				msmtUrl = ""
// 			}
// 		},
// 	)

// 	notifier.Box.AddSubscriber(subscriber)
// 	ctrl.AddOnRecv(func(e devmanager.Event) {
// 		// call when msg recv

// 		cid := config.Params["cid"].(string)
// 		b, err := json.Marshal(map[string]interface{}{
// 			"did":    did,
// 			"cid":    cid,
// 			"status": e.Params()["body"],
// 		})
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		if len(msmtUrl) != 0 {
// 			req, err := http.NewRequest(
// 				"PUT",
// 				msmtUrl,
// 				bytes.NewReader(b),
// 			)

// 			if err != nil {
// 				log.Println(err)
// 				return
// 			}

// 			_, err = http.DefaultClient.Do(req)
// 			if err != nil {
// 				log.Println(err)
// 				return
// 			}
// 		}

// 	})

// 	ctrl.AddOnClose(func(dname, did string, ctrl devmanager.DeviceControllerI) error {
// 		// call when msg recv
// 		cid := config.Params["cid"].(string)
// 		// send request to server for deletion of device
// 		bodyB, err := json.Marshal(map[string]interface{}{
// 			"cid": cid,
// 			"did": did,
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		req, err := http.NewRequest(
// 			"DELETE",
// 			fmt.Sprintf("http://%s/api/v1/devs", config.Params["serverAddr"].(string)),
// 			bytes.NewReader(bodyB),
// 		)
// 		if err != nil {
// 			return err
// 		}

// 		resp, err := http.DefaultClient.Do(req)
// 		if err != nil {
// 			return err
// 		}

// 		if resp.StatusCode > 300 {
// 			b, err := ioutil.ReadAll(resp.Body)
// 			if err != nil {
// 				return err
// 			}

// 			log.Println(string(b))
// 		}
// 		// log.Println(string(b))

// 		// delete controller from cache
// 		cache.RemoveDeviceController(dname)
// 		cache.RemoveDeviceFromSvc(did)
// 		notifier.Box.RemoveSubscriber(subscriber)
// 		return nil
// 	})

// 	return ctrl
// }

// func manageSubscribe() context.CancelFunc {

// 	ctx, cancel := context.WithCancel(context.Background())

// 	go commonutils.Subscribe(
// 		ctx,
// 		"/push/v1/", notifier.SubtokenStatusChanged,
// 		config.Params["cid"].(string),
// 		func(payload []byte) {
// 			// fmt.Println("SUBTOKENSTATUSCHANGED: ", string(payload))
// 			event := map[string]interface{}{}
// 			err := json.Unmarshal(payload, &event)
// 			if err != nil {
// 				log.Println(err)
// 				return
// 			}

// 			key, ok := event["key"].(string)
// 			if !ok {
// 				return
// 			}

// 			if key == "service" {
// 				sname, ok := event["value"].(string)
// 				if !ok {
// 					return
// 				}

// 				sid, err := querySvcID(sname)
// 				if err != nil {
// 					log.Println(err.Error())
// 					return
// 				}
// 				// add service id to cache and subscribe service
// 				cache.AddSvcId(sname, sid)

// 				// register devices to service if service is running
// 				cache.GetSvcList()

// 				devList := cache.GetDevicesOnSvc(sname)
// 				for _, e := range devList {
// 					err = registerDeviceToService(sname, e)
// 					if err != nil {
// 						log.Println(err)
// 					}
// 				}
// 			}
// 		},
// 		func() {
// 			log.Println("connection with edge server is disconnected")
// 			log.Println("please restart after edge server is restarted")
// 			interrupt <- os.Interrupt
// 		},
// 	)

// 	return cancel
// }

// func devManagerTest() {
// 	var line string
// 	reader := bufio.NewReader(os.Stdin)
// 	for {
// 		line, _ = reader.ReadString('\n')
// 		tkns := strings.Split(line, " ")
// 		if tkns[0] == "exit" {
// 			return
// 		}

// 		if tkns[0] == "fan" {
// 			ctrl, err := cache.GetDeviceController("DEVICE-A-UUID")
// 			if err != nil {
// 				panic(err)
// 			}

// 			parameter := false
// 			if tkns[1] == "on\n" {
// 				parameter = true
// 			}

// 			ctrl.Sync(map[string]interface{}{
// 				"fan": parameter,
// 			})
// 		} else if tkns[0] == "lamp" {
// 			ctrl, err := cache.GetDeviceController("DEVICE-A-UUID")
// 			if err != nil {
// 				panic(err)
// 			}

// 			parameter := tkns[1][:2] == "on"
// 			fmt.Println("parameter: ", parameter)
// 			ctrl.Sync(map[string]interface{}{
// 				"lamp": parameter,
// 			})
// 		}
// 	}
// }

var interrupt chan os.Signal

func waitInterrupt() {
	// waiting interrupt
	interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	log.Println("receive interrupt")
}

func main() {
	cfg := flag.Bool("init", false, "create initial config file")
	flag.Parse()

	if *cfg {
		config.CreateInitFile()
		return
	}

	if _, err := os.Stat("./config.properties"); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		fmt.Println("config file doesn't exist")
		fmt.Println("please add -init option to create config file")
		return
	}
	config.LoadConfig()

	statmgmt.Bootup()

	if statmgmt.Status() == statmgmt.STATUS_DISCONNECTED {
		err := statmgmt.Connect()
		if err != nil {
			panic(err)
		}
	}

	// var cancel context.CancelFunc = nil
	// defer func() {
	// 	if cancel != nil {
	// 		cancel()
	// 	}
	// }()
	// if strings.Compare(config.Params["mode"].(string), string(config.MANAGEDBYEDGE)) == 0 {
	// 	cancel = manageSubscribe()
	// }

	// deviceManagerSetup()
	// go devManagerTest()

	go router.NewRouter().Run(config.Params["bind"].(string))

	waitInterrupt()
	// do something before program exit
	// websocket close

}
