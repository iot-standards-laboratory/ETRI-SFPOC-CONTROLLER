package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"etri-sfpoc-controller/common"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/devmanager"
	"etri-sfpoc-controller/model"
	"etri-sfpoc-controller/model/cache"
	"etri-sfpoc-controller/notifier"
	"etri-sfpoc-controller/router"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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
		fmt.Sprintf("http://%s/%s", config.Params["serverAddr"], "api/v1/ctrls"),
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
func querySvcID(sname string) (string, error) {

	// get svc id
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://%s/%s", config.Params["serverAddr"], "api/v1/svcs"),
		nil,
	)

	if err != nil {
		return "", err
	}
	req.Header.Add("sname", sname)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	} else if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", errors.New(string(b))
	}

	b, err := ioutil.ReadAll(resp.Body)
	return string(b), err
}

func deviceManagerSetup() {
	devmanager.AddOnDiscovered(func(port io.ReadWriter, sname, dname string) error {
		defer log.Println("exit onDiscovered()")
		// do register procedure
		//check already registered device

		var did string
		var err error
		did, err = model.DefaultDB.GetDeviceID(dname)
		if err != nil {
			did, err = devmanager.RegisterDevice(map[string]interface{}{
				"sname": sname,
				"dname": dname,
			})
			if err != nil {
				log.Println(err.Error())
				return err
			}

			cid := config.Params["cid"].(string)
			model.DefaultDB.AddDevice(&model.Device{
				DID:   did,
				DName: dname,
				SName: sname,
				CID:   cid,
			})
		} else {
			log.Println("device is already registered")
		}

		// send request to server for registration of device
		okchan := make(chan error)
		defer close(okchan)
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		_, err = port.Write([]byte(`{"code": 1, "token": "initial", "mode": 1}\n`))
		if err != nil {
			return err
		}

		go func() {
			reader := bufio.NewReader(port)
			rcvMsg := map[string]interface{}{}

			for {
				line, _, err := reader.ReadLine()
				if err != nil {
					log.Println("okchan : ", err)
					okchan <- err
					return
				}

				err = json.Unmarshal(line, &rcvMsg)
				if err != nil {
					continue
				}

				fmt.Println("code: ", rcvMsg["code"])
				if rcvMsg["code"].(float64)-1.0 < 0.0001 {
					// register success
					ctrl := makeDeviceController(port, did, dname)
					cache.AddDeviceController(dname, ctrl)
					cache.AddSvc(did, dname, sname)
					sid, err := querySvcID(sname)
					if err == nil {
						log.Println("add sevice ", sname, "'s id : ", sid)
						cache.AddSvcId(sname, sid)
					} else {
						log.Println(err)
					}

					ctrl.Run()
					okchan <- nil
					return
				}
			}
		}()

		for {
			select {
			case <-ticker.C:
				fmt.Println("retransmission command to change mode as timeout")
				_, err = port.Write([]byte(`{"code": 1, "token": "initial", "mode": 1}\n`))
				if err != nil {
					return err
				}
			case err := <-okchan:
				if err == nil {
					return nil
				} else {
					log.Println(err)
				}
			}
		}
	})

	go devmanager.Watch()
}

func makeDeviceController(port io.ReadWriter, did, dname string) devmanager.DeviceControllerI {

	// model.AddDeviceController 에서 등록된 디바이스 목록에 해당 디바이스 추가할 것!!
	ctrl := devmanager.NewDeviceController(port, dname, did)

	ctrl.AddOnRecv(func(e devmanager.Event) {
		// call when msg recv
		fmt.Println(e)
	})

	ctrl.AddOnClose(func(dname, did string, ctrl devmanager.DeviceControllerI) error {
		// call when msg recv

		// send request to server for deletion of device
		// bodyB, err := json.Marshal(map[string]interface{}{
		// 	"dname": dname,
		// 	"did":   did,
		// })
		// if err != nil {
		// 	return err
		// }

		// req, err := http.NewRequest(
		// 	"DELETE",
		// 	fmt.Sprintf("http://%s/api/v1/devs", config.Params["serverAddr"].(string)),
		// 	bytes.NewReader(bodyB),
		// )
		// if err != nil {
		// 	return err
		// }

		// resp, err := http.DefaultClient.Do(req)
		// if err != nil {
		// 	return err
		// }

		// b, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	return err
		// }

		// log.Println(string(b))

		// delete controller from cache
		cache.RemoveDeviceController(dname)
		cache.RemoveDeviceFromSvc(did)
		return nil
	})

	return ctrl
}

func manageSubscribe() {

	go common.Subscribe("/push/v1/", notifier.SubtokenStatusChanged, func(payload []byte) {
		// fmt.Println("SUBTOKENSTATUSCHANGED: ", string(payload))
		event := map[string]interface{}{}
		err := json.Unmarshal(payload, &event)
		if err != nil {
			log.Println(err)
			return
		}

		key, ok := event["key"].(string)
		if !ok {
			return
		}

		if key == "service" {
			sname, ok := event["value"].(string)
			if !ok {
				return
			}

			sid, err := querySvcID(sname)
			if err != nil {
				log.Println(err.Error())
				return
			}
			log.Println("add", sname, "id :", sid)
			cache.AddSvcId(sname, sid)
		}
	})
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
	cid := config.Params["cid"].(string)
	if cid == "blank" {
		var err error
		cid, err = register()
		if err != nil {
			panic(err)
		}

		config.Set("cid", cid)
	}

	manageSubscribe()
	deviceManagerSetup()
	go devManagerTest()
	router.NewRouter().Run(config.Params["bind"].(string))
}

func devManagerTest() {
	var line string
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _ = reader.ReadString('\n')
		tkns := strings.Split(line, " ")
		if tkns[0] == "exit" {
			return
		}

		if tkns[0] == "fan" {
			ctrl, err := cache.GetDeviceController("DEVICE-A-UUID")
			if err != nil {
				panic(err)
			}

			parameter := false
			if tkns[1] == "on\n" {
				parameter = true
			}

			ctrl.Sync(map[string]interface{}{
				"fan": parameter,
			})
		} else if tkns[0] == "lamp" {
			ctrl, err := cache.GetDeviceController("DEVICE-A-UUID")
			if err != nil {
				panic(err)
			}

			parameter := tkns[1][:2] == "on"
			fmt.Println("parameter: ", parameter)
			ctrl.Sync(map[string]interface{}{
				"lamp": parameter,
			})
		}
	}
}
