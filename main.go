package main

import (
	"encoding/json"
	"errors"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/devmanager"
	"etri-sfpoc-controller/model/cachestorage"
	"etri-sfpoc-controller/mqtthandler"
	"etri-sfpoc-controller/router"
	"etri-sfpoc-controller/statmgmt"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
)

// func makeDeviceController(port io.ReadWriter, did, dname, sname string) devmanager.DeviceControllerI {
// 	fmt.Println("makeDeviceController()")
// 	// model.AddDeviceController 에서 등록된 디바이스 목록에 해당 디바이스 추가할 것!!
// 	ctrl := devmanager.NewDeviceController(port, dname, did)

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
	flag.Parse()
	if _, err := os.Stat("./config.properties"); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		fmt.Println("config file doesn't exist")
		config.CreateInitFile()
	}
	config.LoadConfig()

	bindAddr, ok := config.Params["bind"].(string)
	if !ok {
		bindAddr = ":4000"
	}

	go router.NewRouter().Run(bindAddr)

	statmgmt.Bootup()

	if statmgmt.Status() == statmgmt.STATUS_DISCONNECTED {
		err := statmgmt.Connect()
		if err != nil {
			panic(err)
		}
	}

	devmanager.AddOnConnected(func(port string) {
		ctrl, err := devmanager.DiscoverController(port)
		if err != nil {
			log.Println(err)
			return
		}
		ctrl.AddOnUpdate(func(e interface{}) {
			// upload data to mqtt
			id, ok := config.Params["id"].(string)
			if !ok {
				return
			}

			_ = id
			obj := e.(map[string]interface{})
			msg := obj["msg"].(string)
			msgObj := map[string]interface{}{}
			err := json.Unmarshal([]byte(msg), &msgObj)
			if err != nil {
				return
			}

			body := msgObj["body"].(map[string]interface{})
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return
			}

			err = mqtthandler.Publish(fmt.Sprintf("%s/%d", id, ctrl.Key()), bodyBytes)
			if err != nil {
				log.Println("mqtt is disconnected!!")
				log.Println(err)
				return
			}
		})

		ctrl.AddOnError(func(e error) {
			if e.Error() == "EOF" {
				log.Println("EOF error!!")
				ctrl.Close()
			}
		})

		ctrl.AddOnClose(func(key uint64) {
			cachestorage.RemoveDeviceController(key)
		})

		go ctrl.Run()
		cachestorage.AddDeviceController(ctrl)
	})

	go devmanager.Watch()
	// go devManagerTest()

	waitInterrupt()
	// do something before program exit
	// websocket close

}
