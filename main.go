package main

import (
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
	"time"
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

	devmanager.AddOnConnected(func(ctrl devmanager.DeviceControllerI) {
		ctrl.AddOnError(func(e error) {
			if e.Error() == "EOF" {
				log.Println("EOF error!!")
			}
		})

		ctrl.AddOnClose(func(key uint64) {
			cachestorage.RemoveDeviceController(key)
		})

		// add controller to cache and register to edge
		cachestorage.AddDeviceController(ctrl)
		// query controller status on init
		mqtthandler.Subscribe(fmt.Sprintf("%s/%d/post", config.Params["id"], ctrl.Key()))
		mqtthandler.Subscribe(fmt.Sprintf("%s/%d/get", config.Params["id"], ctrl.Key()))
	})

	go devmanager.Watch()

	go func() {
		ticker := time.NewTicker(time.Second * 5)
		keyTemplate := fmt.Sprintf("%s/%%d/content/sensor", config.Params["id"])

		for range ticker.C {
			ctrls := cachestorage.GetDeviceControllers()
			for _, ctrl := range ctrls {
				_, b, err := ctrl.Do(1, []byte("sensor"))
				if err != nil {
					log.Println("error occur: ", err)
					if ctrl.ResetBuffer() != nil {
						ctrl.Close()
						return
					}
					return
				}

				mqtthandler.Publish(fmt.Sprintf(keyTemplate, ctrl.Key()), b)
			}
		}
	}()

	waitInterrupt()
	// do something before program exit
	// websocket close

}
